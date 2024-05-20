package main

import (
	"context"
	"database/sql"
	_ "embed"
	"go-react-embed/api"
	"go-react-embed/frontend"
	"go-react-embed/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main () {
	// load .env file
	err := initAndLoadEnv() // change it to initEnv
	if err != nil {
		log.Fatal("Problem Loading .env\n", err)
	}

	// create DB and Tables and initialize Globals
	initDatabaseModels()
	
	// create echo app
	e := echo.New()
	// middlewares
	e.Use(loggingMiddleware)
	// e.Use(middleware.BodyDump(bodyDump))
	e.Pre(middleware.RemoveTrailingSlash())
	// CORS
	useCORSMiddleware(e)

	// registerign bachend routes routes
	e.GET("/api", root)
	api.RegisterPokemonsHandlers(e.Group("/api"))
	api.RegisterUsersHandlers(e.Group("/api"))
	
	// register react static pages build from react tanstack router
	frontend.RegisterHandlers(e)

	// open app url
	url := os.Getenv("APP_URL")+":"+os.Getenv("APP_PORT")
	err = openBrowser(url)
	if err != nil {
		log.Fatal("Problem Openning the browser\n", err)
	}

	// start server 
	e.Logger.Fatal(e.Start(":"+os.Getenv("APP_PORT")))
}

//go:embed schema/schema.sql
var ddl string

func initDatabaseModels() {
	// connect to database
	databaseFile := "./database.db"
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

// Custom logging middleware
func loggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		end := time.Now()
		latency := end.Sub(start)
		log.Printf("[%s] %s %s %v\n", c.Request().Method, c.Request().URL.Path, c.Request().Proto, latency)
		return err
	}
}

// func bodyDump(c echo.Context, reqBody, resBody []byte) {
// 	fmt.Println("reqBody::", string(reqBody))
// 	fmt.Println("resBody::", string(resBody))
// }

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

func root(ctx echo.Context) error {
	// Defining data
	data := map[string]string{
		"data": "Hello, Gophers.",
	}
	return ctx.JSON(http.StatusOK, data)
}
