package api

import (
	"fmt"
	"go-react-embed/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterUsersHandlers(e *echo.Group) {
	e.GET("/users", getUsersHandler, AuthenticatedMiddleware)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/name/:name", getUserByNameHandler)
	e.PUT("/users", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
}

func getUsersHandler(c echo.Context) error {
	ccu := c.(*CustomContextUser)
	fmt.Println("-ccu:", ccu)
	fmt.Println("-user:", ccu.User)
	
	if ccu.User.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
	}


	users, err := models.QUERIES.ListUsers(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"count": len(users),
		"result": users,
	})
}

func getUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id param must be integer")
	}
	user, err := models.QUERIES.GetUser(models.CTX, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result": user,
	})
}

func getUserByNameHandler(c echo.Context) error {
	name := c.Param("name")
	user, err := models.QUERIES.GetUserByName(models.CTX, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"result": user,
	})
}

func updateUserHandler(c echo.Context) error {
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

	user, err := models.QUERIES.UpdateUser(models.CTX, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update, " + err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "updated successfully",
		"result": user,
	})
}

func deleteUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	err = models.QUERIES.DeleteUser(models.CTX, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to delete " + err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "deleted successfully",
	})
}
