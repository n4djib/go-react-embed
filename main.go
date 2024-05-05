package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main () {
	e := echo.New()
	
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8080", 
			"http://localhost:8081",
		},
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	
	e.GET("/api", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, Gophers....")
	})
	
	// e.GET("*", func(ctx echo.Context) error {
	// 	return ctx.String(http.StatusOK, "catch all routes")
	// })

	// register static (/) page build from react
	// frontend.RegisterHandlers(e)

	// open app url
	// url := "http://localhost:8080/"
    // if err := openURL(url); err != nil {
    //     panic(err)
    // }

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}

// func openURL(url string) error {
//     var cmd *exec.Cmd

//     switch runtime.GOOS {
//     case "windows":
//         cmd = exec.Command("cmd", "/c", "start", url)
//     case "darwin":
//         cmd = exec.Command("open", url)
//     default:
//         cmd = exec.Command("xdg-open", url)
//     }

//     return cmd.Start()
// }
