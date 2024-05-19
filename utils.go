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

func initAndLoadEnv() error {
	err := createEnvFile()
	if err != nil {
		return err
	}

	err = godotenv.Load(".env")
	return err
}

func createEnvFile() error {
	// check if file exist
	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		// create .env file
		f, err := os.Create(".env")
		if err != nil {
			return err
		}
		defer f.Close()

		fmt.Println("Creating .env File...")

		ENV_VARIABLES := `
APP_URL="http://localhost"
APP_PORT="8080"
DEV_PORT="8081"
DATABASE="./database.db"
`
		_, err = f.WriteString(ENV_VARIABLES)
		if err != nil {
			return err
		}
	}
	return nil
}

func openBrowser(url string) error {
	// grab flag
	air_flag := flag.Bool("air", false, "detect if run by AIR")
	flag.Parse()
	air := bool(*air_flag)

	// open app url
	if !air {
		if err := openURL(url); err != nil {
			return err
		}
	}
	return nil
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
