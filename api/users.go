package api

import (
	"go-react-embed/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterUsersHandlers(e *echo.Group) {
	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/name/:name", getUserByNameHandler)
	e.POST("/users", createUserHandler)
	e.PUT("/users", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	e.PUT("/users/active-state", updateUserActiveStateHandler)
}

func getUsersHandler(c echo.Context) error {
	users, err := models.QUERIES.ListUsers(models.CTX)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := DataList{
		Count: len(users),
		Data:  users,
	}
	return c.JSON(http.StatusOK, data)
}

func getUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	user, err := models.QUERIES.GetUser(models.CTX, id)
	// user, err := models.QUERIES.GetUserWithPassword(models.CTX, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := Data{
		Data: user,
	}
	return c.JSON(http.StatusOK, data)
}

func getUserByNameHandler(c echo.Context) error {
	name := c.Param("name")
	user, err := models.QUERIES.GetUserByName(models.CTX, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := Data{
		Data: user,
	}
	return c.JSON(http.StatusOK, data)
}

func createUserHandler(c echo.Context) error {
	var body models.CreateUserParams
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	user, err := models.QUERIES.CreateUser(models.CTX, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "created successfully",
		Data:    user,
	})
}

func updateUserHandler(c echo.Context) error {
	var body models.UpdateUserParams
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	user, err := models.QUERIES.UpdateUser(models.CTX, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "updated successfully",
		Data:    user,
	})
}

func updateUserActiveStateHandler(c echo.Context) error {
	var body models.UpdateUserActiveStateParams
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	user, err := models.QUERIES.UpdateUserActiveState(models.CTX, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "updated status successfully",
		Data:    user,
	})
}

func deleteUserHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	err = models.QUERIES.DeleteUser(models.CTX, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "deleted successfully",
	})
}
