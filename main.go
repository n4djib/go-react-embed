package main

import (
	"flag"
	"go-react-embed/frontend"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Data struct {
	Msg string `json:"msg"`
}

// TODO put urls in ENV

func main () {
	e := echo.New()
	
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8080", 
			"http://localhost:8081",
		},
		// AllowMethods: []string{
		// 	echo.GET, echo.PUT, echo.POST, echo.DELETE
		// },
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	
	e.GET("/api", func(ctx echo.Context) error {
		// return ctx.String(http.StatusOK, "Hello, Gophers....")
		d := &Data{ Msg: "Hello, Gophers."}
		return ctx.JSON(http.StatusOK, d)
	})
	
	// register static (/) page build from react
	frontend.RegisterHandlers(e)

	// open app url
	openBrowser()

	// start server 
	e.Logger.Fatal(e.Start(":8080"))
}

func openBrowser() {
	// grab flag
	air_flag := flag.Bool("air", false, "detect if run by AIR")
	flag.Parse()
	air := bool(*air_flag)
	// fmt.Println(air)

	// open app url
	if !air {
		url := "http://localhost:8080/"
		if err := openURL(url); err != nil {
			panic(err)
		}
	}
}

func openURL(url string) error {
    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "start", url)
    case "darwin":
        cmd = exec.Command("open", url)
    default:
        cmd = exec.Command("xdg-open", url)
    }

    return cmd.Start()
}
