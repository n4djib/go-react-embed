package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
)

func createAndLoadEnv() error {
	// check if file exist
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		// create .env file
		f, err := os.Create(".env")
		check(err)
		defer f.Close()

		fmt.Println("creating .env")

		_, err1 := f.WriteString(`APP_URL="http://localhost"
APP_PORT="8080"
DEV_PORT="8081"`)
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
