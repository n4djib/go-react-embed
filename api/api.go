package api

import (
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
	data := map[string]db.Pokemon{
		"data": pokemon,
	}
	return ctx.JSON(http.StatusOK, data)
}

type PokemonList struct {
	Count int          `json:"count"`
	Data  []db.Pokemon `json:"data"`
}

func pokemonsList(ctx echo.Context) error {
	pokemons, err := db.GetPokemons()
	check(err)

	list := PokemonList{
		Count: len(pokemons),
		Data: pokemons,
	}
	return ctx.JSON(http.StatusOK, list)
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
