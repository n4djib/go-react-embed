package api

import (
	"go-react-embed/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUsersHandlers(e *echo.Group) {
	e.GET("/users", getUsersHandler)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/name/:name", getUserByNameHandler)
	// e.POST("/users", createUserHandler)
	e.PUT("/users", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	e.PUT("/users/active-state", updateUserActiveStateHandler)

	
	e.POST("/auth/signup", signup)
	e.POST("/auth/signin", signin)
}

func signup(c echo.Context) error {
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
	// check user name doesn't exist
	foundUser, err := models.QUERIES.GetUserByName(models.CTX, body.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	if foundUser.ID != 0 {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "Duplicate user name",
		})
	}

	// check password strength

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "Failed to hash password",
		})
	}
	body.Password = string(hash)

	// insert it
	user, err := models.QUERIES.CreateUser(models.CTX, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "Signed up successfully",
		Data:    user,
	})
}


func signin(c echo.Context) error {
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

	// check user name exist
	user, err := models.QUERIES.GetUserByNameWithPassword(models.CTX, body.Name)
	if err != nil || user.ID == 0 {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "user name or password are incorrect 1",
		})
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "user name or password are incorrect 2",
		})
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	// sign the token and encode a a string
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "Failed to create token\n "+ err.Error(),
		})
	}

	user.Password = "[HIDDEN]"
	return c.JSON(http.StatusOK, echo.Map{
		"Message": "Sign in successfully",
		"Data": user,  // we should return JWT
		"token": signedToken,
	})
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

// func createUserHandler(c echo.Context) error {
// 	var body models.CreateUserParams
// 	if err := c.Bind(&body); err != nil {
// 		return c.JSON(http.StatusBadRequest, Error{
// 			Error: err.Error(),
// 		})
// 	}
// 	// Validate the data
// 	if err := validate.Struct(body); err != nil {
// 		return c.JSON(http.StatusBadRequest, Error{
// 			Error: err.Error(),
// 		})
// 	}
// 	user, err := models.QUERIES.CreateUser(models.CTX, body)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, Error{
// 			Error: err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusOK, Status{
// 		Message: "created successfully",
// 		Data:    user,
// 	})
// }

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
