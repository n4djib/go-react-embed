package main

import (
	"errors"
	"flag"
	"fmt"
	"go-react-embed/frontend"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Data struct {
	Msg string `json:"msg"`
}

func main () {
	// load .env file
	err := createAndLoadEnv()
	check(err)
	
	// create echo app
	e := echo.New()
	
	// CORS
	useCORSMiddleware(e)
	
	e.GET("/api", func(ctx echo.Context) error {
		// return ctx.String(http.StatusOK, "Hello, Gophers....")
		d := &Data{ Msg: "Hello, Gophers."}
		return ctx.JSON(http.StatusOK, d)
	})
	
	// register static (/) page build from react
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

func createAndLoadEnv() error {
	// check if file exist
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		// create .env file
		f, err := os.Create(".env");
    	check(err)
		defer f.Close()

		fmt.Println("creating .env")

		_, err1 := f.WriteString(`APP_URL="http://localhost"
APP_PORT="8080"
DEV_PORT="8081"`);
    	check(err1)
	}

	err := godotenv.Load(".env")
	return err
}

func openBrowser(url string) {
	// grab flag
	air_flag := flag.Bool("air", false, "detect if run by AIR")
	flag.Parse()
	air := bool(*air_flag)
	// fmt.Println(air)

	// open app url
	if !air {
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

func check(e error) {
    if e != nil {
        panic(e)
    }
}