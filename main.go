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

	_ "go-react-embed/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// TODO add swagger init to build

// @title GO-REACT-EMBED API
// @version 1.0
// @description this is the API for the backend.
// @termsOfService http://swagger.io/terms/
func main() {
	start := time.Now()

	// load .env file
	initAndLoadEnv()

	// create DB and Tables and initialize Globals
	go initDatabaseModels()

	// setting up RBAC
	go func () {
		rbacAuth, err := setupRBAC()
		if err != nil {
			log.Fatal("+++ Error in Getting RBAC data, ", err)
		}
		// put into gloabal variable
		RBAC = rbacAuth
	}()

	// create echo app
	e := echo.New()
	// adding middlewares
	go addingMiddlewares(e)
	// registering api backend routes
	go registeringApiRoutes(e)
	// register react static pages build from react tanstack router
	go frontend.RegisterHandlers(e)

	// open browser to APP url https://localhost:8080
	go openBrowser()

	// check if file "server.crt", "server.key" exist
	SERVER_CRT, SERVER_KEY := os.Getenv("SERVER_CRT"), os.Getenv("SERVER_KEY")
	go checkSSLFilesExist(SERVER_CRT, SERVER_KEY)

	// hide the banner
	if os.Getenv("HIDE_BANNER") == "true" {
		e.HideBanner = true
	}
	// execution duration
	fmt.Println("- duration:", time.Since(start))
	
	// start server
	APP_PORT := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.StartTLS(":"+APP_PORT, SERVER_CRT, SERVER_KEY))
}

func addingMiddlewares(e *echo.Echo) {
	e.Use(api.WhoamiMiddleware)
	e.Use(loggingMiddleware)
	e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	useCORSMiddleware(e)
	// show it in DEV & SHOW(env)
	// if os.Getenv("MODE") == "DEV" {
	// 	e.Use(middleware.BodyDump(bodyDump))
	// }
}

func registeringApiRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/ping", pong)
	api.RegisterAuthsHandlers(e.Group("/api"))
	api.RegisterPokemonsHandlers(e.Group("/api", api.AuthenticatedMiddleware))
	api.RegisterUsersHandlers(e.Group("/api"))
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
	// ccu := c.(*api.CustomContextUser)
	ctxUser := api.GetUserFromContext(c)
	if ctxUser.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
	}
	
	user := rbac.Map{
		"id": ctxUser.ID, 
		"name": ctxUser.Name, 
		"roles": ctxUser.Roles,
	}
	resource := rbac.Map{"id": 3, "title": "tutorial", "owner": 3, "list": []int{1, 2, 3, 4, 5, 6}}

	start5 := time.Now()
	allowed, err := RBAC.IsAllowed(user, resource, "edit_user")
	if err != nil {
		log.Fatal("++++ error: ", err.Error())
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
