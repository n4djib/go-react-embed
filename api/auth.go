package api

import (
	"fmt"
	"go-react-embed/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAuthsHandlers(e *echo.Group) {
	e.POST("/auth/signup", signup)
	e.POST("/auth/signin", signin)
	e.GET("/auth/signout", signout)
	e.GET("/auth/whoami", whoami)
	e.PUT("/auth/active-state", updateUserActiveStateHandler)
	e.GET("/auth/check-name/:name", checkUsername)
	
	e.GET("/auth/get-rbac", GetRBAC)

	// TODO forgoten password
	// TODO change password
	// TODO check password strength
	// TODO rate limit login tries
	// TODO limit signup tries
}

func GetRBAC(c echo.Context) error {
	roles, err := models.QUERIES.GetRoles(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	permissions, err := models.QUERIES.GetPermissions(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	roleParents, err := models.QUERIES.GetRoleParents(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	permissionParents, err := models.QUERIES.GetPermissionParents(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	rolePermissions, err := models.QUERIES.GetRolePermissions(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fmt.Println("-roles: ")
	for i:=0; i<len(roles); i++ {
		fmt.Println(" ", roles[i])
	}
	fmt.Println("")
	fmt.Println("-permissions: ")
	for i:=0; i<len(permissions); i++ {
		fmt.Println(" ", permissions[i])
	}
	fmt.Println("")
	
	fmt.Println("-roleParents: ")
	for i:=0; i<len(roleParents); i++ {
		fmt.Println(" ", roleParents[i])
	}
	fmt.Println("")
	fmt.Println("-permissionParents: ")
	for i:=0; i<len(permissionParents); i++ {
		fmt.Println(" ", permissionParents[i])
	}
	fmt.Println("")

	fmt.Println("-rolePermissions: ")
	for i:=0; i<len(rolePermissions); i++ {
		fmt.Println(" ", rolePermissions[i])
	}
	fmt.Println("")

	return c.JSON(http.StatusOK, echo.Map{
		"roles": roles,
		"permissions": permissions,
		"roleParents": roleParents,
		"permissionParents": permissionParents,
		"rolePermissions": rolePermissions,
	})
}

type ContextUser struct {
	ID int64
	Name string
	IsActive *bool
	LoggedAt *time.Time
	Roles []string
}

type CustomContextUser struct {
	echo.Context
	User ContextUser
}

// middleware extends the context by adding the authenticated user
func WhoamiMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ccu := &CustomContextUser{c, ContextUser{}}

		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return next(ccu)
		}
		session_uuid := cookie.Value  // uuid
		user, err := models.QUERIES.GetUserBySession(models.CTX, &session_uuid)
		if err != nil {
			return next(ccu)
		}

		// TODO check session expiration

		// get user roles
		roles, _ := models.QUERIES.GetUserRoles(models.CTX, user.ID)
		// if err != nil {
		// 	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to find roles, " + err.Error())
		// }

		cu := &ContextUser{
			user.ID,
			user.Name,
			user.IsActive,
			user.LoggedAt,
			roles,
		}
		
		ccu.User = *cu
		return next(ccu)
	}
}

// AuthenticatedMiddleware checks if token is valid
func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ccu := c.(*CustomContextUser)

		// check if authenticated
		if ccu.User.ID == 0 {
		    return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
		}

		return next(c)
	}
}

func checkUsername(c echo.Context) error {
	name := c.Param("name")
	_, err := models.QUERIES.GetUserByName(models.CTX, name)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "this name is available",
			"exist": false,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "found name",
		"exist": true,
	})
}

func signup(c echo.Context) error {
	var body models.CreateUserParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}
	// check user name doesn't exist
	foundUser, err := models.QUERIES.GetUserByName(models.CTX, body.Name)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if foundUser.ID != 0 {
		return echo.NewHTTPError(http.StatusInternalServerError, "Duplicate user name")
	}

	// TODO check password strength

	// hash password
	hash, err := hashPassword(body.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}
	body.Password = hash
	
	// created at time
	var loggedAt = time.Now()
	body.CreatedAt = &loggedAt

	// insert it
	user, err := models.QUERIES.CreateUser(models.CTX, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Signed up successfully",
		"result": user,
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
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}

	// check user name exist
	user, err := models.QUERIES.GetUserByNameWithPassword(models.CTX, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "user name or password are incorrect")
	}

	// check if user is active
	is_active := *user.IsActive
	if !is_active {
		return echo.NewHTTPError(http.StatusForbidden, "User is not Active")
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
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
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update")
	}

	// 24 hours (1440 minutes)
	COOKIE_EXP_MINUTES := os.Getenv("COOKIE_EXP_MINUTES")
	if len(COOKIE_EXP_MINUTES) < 1 {
		return echo.NewHTTPError(http.StatusInternalServerError, "COOKIE_EXP_MINUTES, can't be found.")
	}
	// convert to int
	expMinutes, err := strconv.ParseInt(COOKIE_EXP_MINUTES, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "COOKIE_EXP_MINUTES, can't be converted.")
	}
	expiration := time.Now().Add(time.Minute * time.Duration(expMinutes))
	// creating cookies
	cookie := createCookie("Authorization", session_uuid, expiration)
	// set cookies
	c.SetCookie(cookie)

	roles, err := models.QUERIES.GetUserRoles(models.CTX, user.ID)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Failed to find user",
			"user": nil,
			"roles": nil,
		})
	}

	user.Password = "[HIDDEN]"
	user.Session = &session_uuid
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Signed in successfully",
		"user": user,
		"roles": roles,
	})
}

func whoami(c echo.Context) error {
	// TODO get it from context not cookie
	// don't forget to check the session expiration
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Authorization Cookie not found",
			"user": nil,
		})
	}
	session_uuid := cookie.Value  // uuid
	user, err := models.QUERIES.GetUserBySession(models.CTX, &session_uuid)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Failed to find user",
			"user": nil,
		})
	}

	// check if session expired
	COOKIE_EXP_MINUTES := os.Getenv("COOKIE_EXP_MINUTES")
	if len(COOKIE_EXP_MINUTES) < 1 {
		return echo.NewHTTPError(http.StatusInternalServerError, "COOKIE_EXP_MINUTES, can't be found.")
	}
	// convert to int
	expMinutes, err := strconv.ParseInt(COOKIE_EXP_MINUTES, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "COOKIE_EXP_MINUTES, can't be converted.")
	}
	expiration := user.LoggedAt.Add(time.Minute * time.Duration(expMinutes))

	if expiration.Before(time.Now()) {
		signout(c)
	}

	roles, err := models.QUERIES.GetUserRoles(models.CTX, user.ID)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Failed to find user",
			"user": nil,
			"roles": nil,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Found whoami",
		"user": user,
		"roles": roles,
	})
}

func signout(c echo.Context) error {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Already Signed out",
		})
	}

	session_uuid := cookie.Value
	user, err := models.QUERIES.GetUserBySession(models.CTX, &session_uuid)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Already Signed out",
		})
	}

	// remove session from users table
	var emptySession models.UpdateUserSessionParams
	emptySession.ID = user.ID
	emptySession.Session = nil
	emptySession.LoggedAt = user.LoggedAt

	err = models.QUERIES.UpdateUserSession(models.CTX, emptySession)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove session, " + err.Error())
	}

	// unset cookie
	deadCookie := createCookie("Authorization", "", time.Unix(0, 0))
	c.SetCookie(deadCookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Sign out successfully",
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
	if os.Getenv("MODE") == "DEV" {
		cookie.SameSite = http.SameSiteNoneMode
	}
	return cookie
}

func updateUserActiveStateHandler(c echo.Context) error {
	var body models.UpdateUserActiveStateParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}
	_, err := models.QUERIES.UpdateUserActiveState(models.CTX, body)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	// user not found
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update, " + err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "updated status successfully",
	})
}
