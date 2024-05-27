package api

import (
	"go-react-embed/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUsersHandlers(e *echo.Group) {
	e.GET("/users", getUsersHandler, AuthenticatedMiddleware)
	e.GET("/users/:id", getUserHandler)
	e.GET("/users/name/:name", getUserByNameHandler)
	e.PUT("/users", updateUserHandler)
	e.DELETE("/users/:id", deleteUserHandler)
	e.PUT("/users/active-state", updateUserActiveStateHandler)

	
	e.POST("/auth/signup", signup)
	e.POST("/auth/signin", signin)
	e.GET("/auth/signout", signout)

	// TODO
	// forgoten password
	// add a check username exist to ue at signup
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

	// TODO check password strength

	// hash password
	hash, err := hashPassword(body.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "Failed to hash password",
		})
	}
	body.Password = hash
	
	// created at time
	var loggedAt = time.Now()
	body.CreatedAt = &loggedAt

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

func hashPassword(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
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
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "user name or password are incorrect 1",
		})
	}

	// check if user is active
	is_active := *user.IsActive
	if !is_active {
		return c.JSON(http.StatusBadRequest, Error{
			Error: "User is not Active",
		})
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "user name or password are incorrect 2",
		})
	}

	var updateSession models.UpdateUserSessionParams
	var loggedAt = time.Now()
	session_uuid := uuid.New().String()

	updateSession.ID = user.ID
	updateSession.Session = &session_uuid
	updateSession.LoggedAt = &loggedAt

	// put session in user table
	err = models.QUERIES.UpdateUserSession(models.CTX, updateSession)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}

	// 24 hours (1440 minutes)
	COOKIE_EXP_MINUTES := os.Getenv("COOKIE_EXP_MINUTES")
	if len(COOKIE_EXP_MINUTES) < 1 {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "COOKIE_EXP_MINUTES, can't be found.",
		})
	}
	// convert to int
	expMinutes, err := strconv.ParseInt(COOKIE_EXP_MINUTES, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "COOKIE_EXP_MINUTES, can't be converted," + err.Error(),
		})
	}
	expiration := time.Now().Add(time.Minute * time.Duration(expMinutes))
	// creating cookies
	cookie := createCookie("Authorization", session_uuid, expiration)
	// set cookies
	c.SetCookie(cookie)

	user.Password = "[HIDDEN]"
	return c.JSON(http.StatusOK, echo.Map{
		"Message": "Sign in successfully",
	})
}

// AuthenticatedMiddleware checks if token is valid
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "Authorization Cookie not found, " + err.Error())
		}
		
		session_uuid := cookie.Value  // uuid
		_, err = models.QUERIES.GetUserBySession(models.CTX, &session_uuid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find user 1, " + err.Error())
		}

		return next(c)
	}
}

func signout(c echo.Context) error {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Already Signed out 1, " + err.Error(),
		})
	}

	session_uuid := cookie.Value
	user, err := models.QUERIES.GetUserBySession(models.CTX, &session_uuid)
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Already Signed out 2, " + err.Error(),
		})
	}

	// remove session from users table
	var emptySession models.UpdateUserSessionParams
	emptySession.ID = user.ID
	emptySession.Session = nil
	emptySession.LoggedAt = user.LoggedAt

	err = models.QUERIES.UpdateUserSession(models.CTX, emptySession)
	if err != nil {
		return c.JSON(http.StatusOK, Error{
			Error: "Failed to remove session, " + err.Error(),
		})
	}

	// unset cookie
	deadCookie := createCookie("Authorization", "", time.Unix(0, 0))
	c.SetCookie(deadCookie)

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
	if timeHoursExpiry == time.Unix(0, 0) {
		cookie.MaxAge = -1
	}
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

	hash, err := hashPassword(body.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: "Failed to hash password",
		})
	}
	body.Password = hash

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
	err := models.QUERIES.UpdateUserActiveState(models.CTX, body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, Status{
		Message: "updated status successfully",
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
