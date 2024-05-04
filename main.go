package main

import (
	"go-react-embed/frontend"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/labstack/echo/v4"
)

func main () {
	e := echo.New()
	
	e.GET("/api", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, Gophers!")
	})

	frontend.RegisterHandlers(e)

	// open app url
	url := "http://localhost:8080/"
    err := openURL(url)
    if err != nil {
        panic(err)
    }

	// start server
	e.Logger.Fatal(e.Start(":8080"))
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
