package api

import (
	"fmt"
	"go-react-embed/db"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Group) {
	e.GET("", root)
	e.GET("/hello", sayHello)
	e.GET("/pokemons", pokemonsList)
	e.GET("/pokemons/:id", pokemon)

	e.POST("/auth/signup", addUser)
	e.GET("/auth/users", usersList)
}

func root(ctx echo.Context) error {
	// Defining data
	data := map[string]string{
		"data": "Hello, Gophers.",
	}
	return ctx.JSON(http.StatusOK, data)
}

func sayHello(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello")
}

func pokemon(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	check(err)

	pokemon, err := db.GetPokemon(id)
	check(err)

	// Defining data
	data := db.Data{
		Data: pokemon,
	}
	return ctx.JSON(http.StatusOK, data)
}

func pokemonsList(ctx echo.Context) error {
	pokemons, err := db.GetPokemons()
	check(err)
	data := db.DataList{
		Count: len(pokemons),
		Data: pokemons,
	}
	return ctx.JSON(http.StatusOK, data)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


/////////////////
 
func usersList(ctx echo.Context) error {
	users, err := db.GetUsers()
	check(err)
	data := db.DataList{
		Count: len(users),
		Data: users,
	}
	return ctx.JSON(http.StatusOK, data)
}

func addUser(ctx echo.Context) error {
	var jsonUser db.User
	err := ctx.Bind(&jsonUser)
	check(err)
	created, err := db.CreateUser(jsonUser)
	fmt.Println(created, " - ", err)
	return ctx.JSON(http.StatusOK, created)
}
