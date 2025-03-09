package api

import (
	"context"
	"database/sql"
	"go-react-embed/models"
	"go-react-embed/rbac"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	CTX     context.Context
	QUERIES *models.Queries
	DB      *sql.DB
	RBAC    rbac.RBAC
}

func (h AuthHandler) RegisterHandlers(e *echo.Group) {
	e.POST("/auth/signup", h.signup)
	e.POST("/auth/signin", h.signin)
	e.GET("/auth/signout", h.signout)
	e.GET("/auth/whoami", h.whoami)
	e.PUT("/auth/active-state", h.updateUserActiveStateHandler)
	e.GET("/auth/check-name/:name", h.checkUsername)
	
	e.GET("/auth/get-rbac", h.GetRBAC)

	// TODO forgoten password
	// TODO change password
	// TODO check password strength
	// TODO rate limit login tries
	// TODO limit signup tries
}

func (h AuthHandler) GetRBAC(c echo.Context) error {
	roles, err := h.QUERIES.GetRoles(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	permissions, err := h.QUERIES.GetPermissions(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	roleParents, err := h.QUERIES.GetRoleParents(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	permissionParents, err := h.QUERIES.GetPermissionParents(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	rolePermissions, err := h.QUERIES.GetRolePermissions(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// fmt.Println("-roles: ")
	// for i:=0; i<len(roles); i++ {
	// 	fmt.Println(" ", roles[i])
	// }
	// fmt.Println("")
	// fmt.Println("-permissions: ")
	// for i:=0; i<len(permissions); i++ {
	// 	fmt.Println(" ", permissions[i])
	// }
	// fmt.Println("")
	
	// fmt.Println("-roleParents: ")
	// for i:=0; i<len(roleParents); i++ {
	// 	fmt.Println(" ", roleParents[i])
	// }
	// fmt.Println("")
	// fmt.Println("-permissionParents: ")
	// for i:=0; i<len(permissionParents); i++ {
	// 	fmt.Println(" ", permissionParents[i])
	// }
	// fmt.Println("")

	// fmt.Println("-rolePermissions: ")
	// for i:=0; i<len(rolePermissions); i++ {
	// 	fmt.Println(" ", rolePermissions[i])
	// }
	// fmt.Println("")

	return c.JSON(http.StatusOK, echo.Map{
		"roles": roles,
		"permissions": permissions,
		"roleParents": roleParents,
		"permissionParents": permissionParents,
		"rolePermissions": rolePermissions,
	})
}

type ContextUser struct {
	ID int64            `json:"id"`
	Name string         `json:"name"`
	IsActive *bool      `json:"is_active"`
	LoggedAt *time.Time `json:"logged_at"`
	Roles []string      `json:"roles"`
}

type CustomContextUser struct {
	echo.Context
	User ContextUser
}

func (cu CustomContextUser) GetUser() ContextUser {
	return cu.User
}

func GetUserFromContext(c echo.Context) ContextUser {
	ccu := c.(*CustomContextUser)
	return ccu.GetUser()
}

// middleware extends the context by adding the authenticated user
func (h AuthHandler) WhoamiMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ccu := &CustomContextUser{c, ContextUser{}}

		cookie, err := c.Cookie("Authorization")
		if err != nil {
			return next(ccu)
		}
		session_uuid := cookie.Value  // uuid
		user, err := h.QUERIES.GetUserBySession(h.CTX, &session_uuid)
		if err != nil {
			return next(ccu)
		}

		// TODO check session expiration

		// get user roles
		roles, _ := h.QUERIES.GetUserRoles(h.CTX, user.ID)
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
		// ccu := c.(*CustomContextUser)
		user := GetUserFromContext(c)

		// check if authenticated
		if user.ID == 0 {
		    return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
		}

		return next(c)
	}
}

func (h AuthHandler) checkUsername(c echo.Context) error {
	name := c.Param("name")
	_, err := h.QUERIES.GetUserByName(h.CTX, name)
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

func (h AuthHandler) signup(c echo.Context) error {
	var body models.CreateUserParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}
	// check user name doesn't exist
	foundUser, err := h.QUERIES.GetUserByName(h.CTX, body.Name)
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
	user, err := h.QUERIES.CreateUser(h.CTX, body)
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

func (h AuthHandler) signin(c echo.Context) error {
	var body models.CreateUserParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}

	// check user name exist
	user, err := h.QUERIES.GetUserByNameWithPassword(h.CTX, body.Name)
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
	err = h.QUERIES.UpdateUserSession(h.CTX, updateSession)
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

	roles, err := h.QUERIES.GetUserRoles(h.CTX, user.ID)
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

func (h AuthHandler) whoami(c echo.Context) error {
	ccu := c.(*CustomContextUser)
	if ccu.User.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
	}

	user := ccu.User

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
		h.signout(c)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Found whoami",
		"user": user,
	})
}

func (h AuthHandler) signout(c echo.Context) error {
	cookie, err := c.Cookie("Authorization")
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Already Signed out",
		})
	}

	session_uuid := cookie.Value
	user, err := h.QUERIES.GetUserBySession(h.CTX, &session_uuid)
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

	err = h.QUERIES.UpdateUserSession(h.CTX, emptySession)
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

func (h AuthHandler) updateUserActiveStateHandler(c echo.Context) error {
	var body models.UpdateUserActiveStateParams
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	}
	// Validate the data
	if err := validate.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to validate, " + err.Error())
	}
	_, err := h.QUERIES.UpdateUserActiveState(h.CTX, body)
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
