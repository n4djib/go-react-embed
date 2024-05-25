package api

import (
	"errors"
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
	e.GET("/users", getUsersHandler, AuthenticatedMiddleware)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/name/:name", getUserByNameHandler)
	// e.POST("/users", createUserHandler)
	e.PUT("/users", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	e.PUT("/users/active-state", updateUserActiveStateHandler)

	
	e.POST("/auth/signup", signup)
	e.POST("/auth/signin", signin)
	e.GET("/auth/validate-token", validateToken) // TODO remove this
	e.GET("/auth/signout", signout)

	// TODO
	// signout
	//    remove session
	//    delete cookie
	// forgoten password
	// edit user & password
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

type UserClaim struct {
    Sub int64 `json:"sub"`
    Exp int64 `json:"exp"`
    jwt.RegisteredClaims
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

	// 
	secret := os.Getenv("JWT_SECRET")
	exp :=  time.Now().Add(time.Hour * 24 * 30).Unix()

	// generate JWT token
	signedToken, err := generateUserSignedToken(user, secret, exp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}

	expiration := time.Now().Add(24 * time.Hour)
	// creating cookies
	cookie := createCookie("Authorization", signedToken, expiration)
	// set cookies
	c.SetCookie(cookie)

	user.Password = "[HIDDEN]"
	return c.JSON(http.StatusOK, echo.Map{
		"Message": "Sign in successfully",
		// "Data": user,  // we should return JWT
		"token": signedToken,
	})
}

func generateUserSignedToken(user models.User, secret string, expiry int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		Sub: user.ID,
		Exp: expiry,
	})
	// sign the token and encode a a string
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("Failed to create token\n "+ err.Error())
	}
	return signedToken, nil
}

func parseSignedToken(signedToken string, secret string) (UserClaim, error) {
	var userClaim UserClaim
	token, err := jwt.ParseWithClaims(signedToken, &userClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return UserClaim{}, err
	}
	// Checking token validity 
	if !token.Valid {
		return UserClaim{}, errors.New("token is not valid!")
	}
	return userClaim, nil
}

func validateToken(c echo.Context) error {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Cookie not found, " + err.Error(),
		})
	}

	secret := os.Getenv("JWT_SECRET")
	signedToken := cookie.Value
	userClam, err := parseSignedToken(signedToken , secret)
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Failed to validate token, " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"Message": "Token is valid",
		"User ID": userClam.Sub,
	})
}

// AuthenticatedMiddleware checks if token is valid
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Cookie not found, " + err.Error())
		}
		secret := os.Getenv("JWT_SECRET")
		signedToken := cookie.Value
		_, err = parseSignedToken(signedToken , secret)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Failed to validate token, " + err.Error())
		}
		return next(c)
	}
}

func signout(c echo.Context) error {
	_, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Already Signed out, " + err.Error(),
		})
	}

	// unset cookie
	c.SetCookie(unsetCookie("Authorization"))

	return c.JSON(http.StatusOK, echo.Map{
		"Message": "Sign out successfully",
	})
}

func createCookie(name string, value string, timeHoursExpiry time.Time) *http.Cookie {
	cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = value
    cookie.Path = "/"
    cookie.Expires = timeHoursExpiry
    cookie.HttpOnly = true
    cookie.Secure = true
	return cookie
}

func unsetCookie(name string) *http.Cookie {
	cookie := new(http.Cookie)
    cookie.Name = name
    cookie.Value = ""
    cookie.Path = "/"
    cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
    cookie.HttpOnly = true
    cookie.Secure = true
	return cookie
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
