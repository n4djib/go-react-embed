package main

import (
	_ "embed"
	"fmt"

	"go-react-embed/api"
	"go-react-embed/frontend"
	"go-react-embed/rbac"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "go-react-embed/docs"
)

// TODO add swag init to build

// @title GO-REACT-EMBED API
// @version 1.0
// @description this is the API for the backend.
// @termsOfService http://swagger.io/terms/
func main() {
	// load .env file
	err := initAndLoadEnv()
	if err != nil {
		log.Fatal("Problem Loading .env\n", err)
	}

	// create DB and Tables and initialize Globals
	initDatabaseModels()

	// rbacAuth, err := setupRBAC()
	rbacAuth, err := setupRBAC()
	if err != nil {
		log.Fatal("+++ Error in Getting RBAC data, ", err)
	}
	RBAC = rbacAuth

	// create echo app
	e := echo.New()
	e.Use(api.CurrentAuthUserMiddleware)
	// middlewares
	e.Use(loggingMiddleware)
	e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	useCORSMiddleware(e)
	// TODO show it in DEV & SHOW(env)
	// if os.Getenv("MODE") == "DEV" {
	// 	e.Use(middleware.BodyDump(bodyDump))
	// }

	// registering bachend routes routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// e.GET("/ping", pong, api.AuthenticatedMiddleware)
	e.GET("/ping", pong)
	api.RegisterPokemonsHandlers(e.Group("/api", api.AuthenticatedMiddleware))
	api.RegisterUsersHandlers(e.Group("/api"))
	api.RegisterAuthsHandlers(e.Group("/api"))

	// register react static pages build from react tanstack router
	frontend.RegisterHandlers(e)

	// open browser to APP url https://localhost:8080
	err = openBrowser()
	if err != nil {
		log.Fatal("Problem Opening the browser\n", err)
	}

	// check if file "server.crt", "server.key" exist
	SERVER_CRT := os.Getenv("SERVER_CRT")
	SERVER_KEY := os.Getenv("SERVER_KEY")
	err = checkSSLFilesExist(SERVER_CRT, SERVER_KEY)
	if err != nil {
		log.Fatal("SSL files not found\n", err)
	}

	// FIXME echo: http: TLS handshake error from [::1]:49955:
	// TODO hide the banner
	// start server
	e.Logger.Fatal(e.StartTLS(":"+os.Getenv("APP_PORT"), SERVER_CRT, SERVER_KEY))
}

// Custom logging middleware
func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		end := time.Now()
		latency := end.Sub(start)
		log.Printf("[%s] %s %s %v\n", c.Request().Method, c.Request().URL.Path, c.Request().Proto, latency)
		return err
	}
}

// func bodyDump(c echo.Context, reqBody, resBody []byte) {
// 	fmt.Println("reqBody::", string(reqBody))
// 	fmt.Println("resBody::", string(resBody))
// }

func useCORSMiddleware(e *echo.Echo) {
	allowOrigins := []string{os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT")}
	if os.Getenv("MODE") == "DEV" {
		allowOrigins = append(allowOrigins, os.Getenv("APP_URL_DEV")+":"+os.Getenv("APP_PORT_DEV"))
	}

	allowMethods := []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}

	corsConfig := middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
		AllowCredentials: true,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
}

func pong(c echo.Context) error {
	// testing authorization

	// ucc := c.(*api.CustomContextUser)
	// fmt.Println("-ucc:", ucc)
	
	user := rbac.Map{
		"id": 5, "name": "nadjib", "age": 4, 
		"roles": []string{
			// "ADMIN", 
			"USER", 
		},
	}

	ressource := rbac.Map{"id": 5, "title": "tutorial", "owner": 5, "list": []int{1, 2, 3, 4, 5, 6}}

	start5 := time.Now()
	allowed, err := RBAC.IsAllowed(user, ressource, "edit_own_user")
	if err != nil {
		log.Fatal("++++ error: ", err.Error())
	}
	fmt.Println("-allowed:", allowed)
	fmt.Println("-duration5:", time.Since(start5))

	// Defining data
	data := map[string]string{
		"message": "Pong!",
	}
	return c.JSON(http.StatusOK, data)
}
