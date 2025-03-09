package api

import (
	"context"
	"database/sql"
	"fmt"
	"go-react-embed/models"
	"go-react-embed/rbac"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	CTX     context.Context
	QUERIES *models.Queries
	DB      *sql.DB
	RBAC    rbac.RBAC
}

func (h UserHandler) RegisterHandlers(e *echo.Group) {
	e.GET("/users", h.getUsersHandler, AuthenticatedMiddleware)
	e.GET("/users/:id", h.getUserHandler)
	e.GET("/users/name/:name", h.getUserByNameHandler)
	e.PUT("/users", h.updateUserHandler)
	e.DELETE("/users/:id", h.deleteUserHandler)
}

func (h UserHandler) getUsersHandler(c echo.Context) error {
	ccu := c.(*CustomContextUser)
	fmt.Println("-ccu:", ccu)
	fmt.Println("-user:", ccu.User)
	
	if ccu.User.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
	}

	users, err := h.QUERIES.ListUsers(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"count": len(users),
		"result": users,
	})
}

func (h UserHandler) getUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id param must be integer")
	}
	user, err := h.QUERIES.GetUser(h.CTX, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result": user,
	})
}

func (h UserHandler) getUserByNameHandler(c echo.Context) error {
	name := c.Param("name")
	user, err := h.QUERIES.GetUserByName(h.CTX, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result": user,
	})
}

func (h UserHandler) updateUserHandler(c echo.Context) error {
	var body models.UpdateUserParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}
	body.Password = hash

	user, err := h.QUERIES.UpdateUser(h.CTX, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update, " + err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "updated successfully",
		"result": user,
	})
}

func (h UserHandler) deleteUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	err = h.QUERIES.DeleteUser(h.CTX, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to delete " + err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "deleted successfully",
	})
}
