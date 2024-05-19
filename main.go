package main

import (
	"go-react-embed/api"
	"go-react-embed/db"
	"go-react-embed/frontend"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main () {
	// load .env file
	err := createAndLoadEnv()
	check(err)

	// create db & tables if not exist
	err = db.CreateDbTables(os.Getenv("DATABASE"))
	check(err)
	defer db.CloseDatabase()
	
	// create echo app
	e := echo.New()
	
	// CORS
	useCORSMiddleware(e)

	// from routes file
	api.RegisterHandlers(e.Group("/api"))
	
	// register react static pages build from react
	frontend.RegisterHandlers(e)

	// open app url
	url := os.Getenv("APP_URL")+":"+os.Getenv("APP_PORT")
	openBrowser(url)

	// start server 
	e.Logger.Fatal(e.Start(":"+os.Getenv("APP_PORT")))
}

func useCORSMiddleware(e *echo.Echo) {
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{
			os.Getenv("APP_URL")+":"+os.Getenv("APP_PORT"), 
			os.Getenv("APP_URL")+":"+os.Getenv("DEV_PORT"),
		},
		AllowMethods: []string{
			echo.GET, echo.PUT, echo.POST, echo.DELETE,
		},
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
}
