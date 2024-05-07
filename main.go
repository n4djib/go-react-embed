package main

import (
	"flag"
	"go-react-embed/frontend"
	"log"
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

// TODO put urls in ENV

func main () {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	
	// create echo app
	e := echo.New()
	
	// CORS
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
	// e.Logger.Fatal(e.Start( ":8080" ))
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
