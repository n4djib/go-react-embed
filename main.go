package main

import (
	_ "embed"
	"fmt"
	"go-react-embed/api"
	"go-react-embed/frontend"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main () {
	// load .env file
	err := initAndLoadEnv()
	if err != nil {
		log.Fatal("Problem Loading .env\n", err)
	}

	// create DB and Tables and initialize Globals
	initDatabaseModels()
	
	// create echo app
	e := echo.New()
	// middlewares
	e.Use(loggingMiddleware)
	e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	useCORSMiddleware(e)
	// if os.Getenv("MODE") == "DEV" {
	// 	e.Use(middleware.BodyDump(bodyDump))
	// }

	// registering bachend routes routes
	e.GET("/api", root)
	api.RegisterPokemonsHandlers(e.Group("/api"))
	api.RegisterUsersHandlers(e.Group("/api"))
	
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

func bodyDump(c echo.Context, reqBody, resBody []byte) {
	fmt.Println("reqBody::", string(reqBody))
	fmt.Println("resBody::", string(resBody))
}

func useCORSMiddleware(e *echo.Echo) {
	allowOrigins := []string{os.Getenv("APP_URL")+":"+os.Getenv("APP_PORT")}
	if os.Getenv("MODE") == "DEV" {
		allowOrigins = append(allowOrigins, os.Getenv("APP_URL_DEV")+":"+os.Getenv("APP_PORT_DEV"))
	}
	
	allowMethods := []string{echo.GET, echo.PUT, echo.POST, echo.DELETE}

	corsConfig := middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: allowMethods,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
}

func root(ctx echo.Context) error {
	// Defining data
	data := map[string]string{
		"data": "Hello, Gophers.",
	}
	return ctx.JSON(http.StatusOK, data)
}
