package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go-react-embed/models"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func initAndLoadEnv() error {
	err := createEnvFile()
	if err != nil {
		return err
	}
	err = godotenv.Load(".env")
	return err
}

func openBrowser() error {
	url := os.Getenv("APP_URL") + ":" + os.Getenv("APP_PORT")
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

//go:embed schema/schema.sql
var ddl string

func initDatabaseModels() {
	// connect to database
	databaseFile := os.Getenv("DATABASE")
	db, err := sql.Open("sqlite3", databaseFile)
	if err != nil {
		log.Fatal("Connection to DB error\n", err)
	}
	// createTables
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal("Table Cretation error\n", err)
	}
	queries := models.New(db)
	// assign to global variables in models package
	models.DB, models.CTX, models.QUERIES = db, ctx, queries
}

func checkSSLFilesExist(SERVER_CRT string, SERVER_KEY string) error {
	if _, err := os.Stat(SERVER_CRT); errors.Is(err, os.ErrNotExist) {
		return errors.New("Can't find file: " + SERVER_CRT)
	}
	if _, err := os.Stat("server.key"); errors.Is(err, os.ErrNotExist) {
		return errors.New("Can't find file: " + SERVER_KEY)
	}
	return nil
}

func createEnvFile() error {
	// check if file exist
	_, err := os.Stat(".env"); 
	if !errors.Is(err, os.ErrNotExist) {
		return nil
	}
	
	// create .env file
	f, err := os.Create(".env")
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Println("Creating .env File...")

	// dev url and port are not needed after build
	ENV_VARIABLES := `
MODE="PRODUCTION"
APP_URL="https://localhost"
APP_PORT="8080"
DATABASE="./go-react-embed.db"
SERVER_CRT="server.crt"
SERVER_KEY="server.key"
`
	_, err = f.WriteString(ENV_VARIABLES)
	return err
}
