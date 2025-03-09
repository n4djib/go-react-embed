package api

import (
	"context"
	"database/sql"
	"fmt"
	"go-react-embed/models"
	"go-react-embed/rbac"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type PingHandler struct {
	CTX     context.Context
	QUERIES *models.Queries
	DB      *sql.DB
	RBAC    rbac.RBAC
}

func (h PingHandler) RegisterHandlers(e *echo.Group) {
	e.GET("/ping", h.pong)
}

func (h PingHandler) pong(c echo.Context) error {
	// testing authorization
	// ccu := c.(*api.CustomContextUser)
	ctxUser := GetUserFromContext(c)
	if ctxUser.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated (checked inside pong)")
	}
	
	user := rbac.Map{
		"id": ctxUser.ID, 
		"name": ctxUser.Name, 
		"roles": ctxUser.Roles,
	}
	resource := rbac.Map{"id": 3, "title": "tutorial", "owner": 3, "list": []int{1, 2, 3, 4, 5, 6}}

	start5 := time.Now()
	allowed, err := h.RBAC.IsAllowed(user, resource, "edit_user")
	if err != nil {
		fmt.Println("++++ error: ", err.Error())
	}
	fmt.Println("-allowed to ping:", allowed)
	fmt.Println("-duration5:", time.Since(start5))

	if !allowed {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authorized")
	}

	// Defining data
	data := map[string]string{
		"message": "Pong!",
	}
	return c.JSON(http.StatusOK, data)
}
