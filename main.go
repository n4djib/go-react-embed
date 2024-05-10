package main

import (
	"go-react-embed/api"
	"go-react-embed/frontend"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main () {
	// load .env file
	err := createAndLoadEnv()
	check(err)
	
	// create echo app
	e := echo.New()
	
	// CORS
	useCORSMiddleware(e)
	
	// e.GET("/api", func(ctx echo.Context) error {
	//  return ctx.String(http.StatusOK, "Hello, Gophers....")
	// })

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
		// AllowMethods: []string{
		// 	echo.GET, echo.PUT, echo.POST, echo.DELETE
		// },
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
}

