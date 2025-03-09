package main

import (
	"fmt"
	"go-react-embed/api"
	"go-react-embed/frontend"
	"log"
	"os"
	"time"

	_ "go-react-embed/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title GO-REACT-EMBED API
// @version 1.0
// @description this is the API for the backend.
// @termsOfService http://swagger.io/terms/
func main() {
	start := time.Now()
	fmt.Println("started at:", start)
	
	// load .env file
	initAndLoadEnv()
	
	// create echo app
	e := echo.New()

	// adding middlewares
	go addingMiddlewares(e)

	// registering api backend routes
	go func () {
		// create DB and Tables and initialize Globals
		db, ctx, queries := initDatabaseModels()

		// setting up RBAC
		rbacAuth, err := setupRBAC(ctx, queries)
		if err != nil {
			// TODO show the error in Cnsole and freeze for the user to see the problem
			// maybe create a printing func
			fmt.Println("+++ Error in Getting RBAC data, ", err)
		}

		// regidtering handlers
		pingHandlers := api.PingHandler{CTX: ctx, QUERIES: queries, RBAC: rbacAuth, DB: db}
		pokemonHandlers := api.PokemonHandler{CTX: ctx, QUERIES: queries, RBAC: rbacAuth, DB: db}
		userHandlers := api.UserHandler{CTX: ctx, QUERIES: queries, RBAC: rbacAuth, DB: db}
		authHandlers := api.AuthHandler{CTX: ctx, QUERIES: queries, RBAC: rbacAuth, DB: db}

		pingHandlers.RegisterHandlers(e.Group("/api"))
		authHandlers.RegisterHandlers(e.Group("/api"))
		// pokemonHandlers.RegisterHandlers(e.Group("/api", api.AuthenticatedMiddleware))
		pokemonHandlers.RegisterHandlers(e.Group("/api"))
		userHandlers.RegisterHandlers(e.Group("/api"))
		e.Use(authHandlers.WhoamiMiddleware)

		// TODO add swagger init to build
		// registering swagger api
		e.GET("/api/swagger/*", echoSwagger.WrapHandler, api.AuthenticatedMiddleware)
		// https://localhost:8080/api/swagger/index.html

		
		// register react static pages build from react tanstack router
		frontend.RegisterHandlers(e)
	}()

	// open browser to APP url https://localhost:8080
	go openBrowser()
	
	// check if file "server.crt", "server.key" exist
	SERVER_CRT, SERVER_KEY := os.Getenv("SERVER_CRT"), os.Getenv("SERVER_KEY")
	// go 
	// FIXME the generated files are not valid !
	checkSSLFilesExist(SERVER_CRT, SERVER_KEY)

	// hide the banner
	if os.Getenv("HIDE_BANNER") == "true" {
		e.HideBanner = true
	}
	// execution duration
	finished := time.Now()
	fmt.Println("finished at:", finished)
	fmt.Println("- Duration:", time.Since(start))
	
	// start server
	APP_PORT := os.Getenv("APP_PORT")
	e.Logger.Fatal(e.StartTLS(":"+APP_PORT, SERVER_CRT, SERVER_KEY))
}

func addingMiddlewares(e *echo.Echo) {
	// e.Use(WhoamiMiddleware)
	e.Use(loggingMiddleware)
	e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	useCORSMiddleware(e)
	// show it in DEV & SHOW(env)
	// if os.Getenv("MODE") == "DEV" {
	// 	e.Use(middleware.BodyDump(bodyDump))
	// }
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

func bodyDump(c echo.Context, reqBody, resBody []byte) {
	fmt.Println("reqBody::", string(reqBody))
	fmt.Println("resBody::", string(resBody))
}

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
