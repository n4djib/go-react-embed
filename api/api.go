package api

import (
	"fmt"
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

type Data struct {
	Msg string `json:"msg"`
}

type Pokemon struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`
}

type PokemonList struct {
	Count int `json:"count"`
	Results []Pokemon `json:"results"`
}

func root(ctx echo.Context) error {
	data := &Data{ Msg: "Hello, Gophers."}
	return ctx.JSON(http.StatusOK, data)
}

func sayHello(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Hello")
}

var pokemon1 = Pokemon{
	ID: 1,
	Name: "spearow", 
	Url: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/21.svg",
	// Url: "https://pokeapi.co/api/v2/pokemon/21/",
}

var pokemon2 = Pokemon{
	ID: 2,
	Name: "fearow", 
	Url: "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/dream-world/22.svg",
	// Url: "https://pokeapi.co/api/v2/pokemon/22/",
}

var pokemons = []Pokemon{
	pokemon1,
	pokemon2,
}

func pokemon(ctx echo.Context) error {
	param := ctx.Param("id")
	id, err := strconv.Atoi(param)
	check(err)
	fmt.Println(id)
	for _, p := range pokemons {
		if p.ID == id {
			return ctx.JSON(http.StatusOK, p)
		}
	}
	return ctx.JSON(http.StatusOK, Pokemon{})
}

func pokemonsList(ctx echo.Context) error {
	list := PokemonList{
		Count: len(pokemons),
		Results: pokemons,
	}
	return ctx.JSON(http.StatusOK, list)
}

// func pokemonsList2(ctx echo.Context) error {
// 	// // fetch pokemon list
// 	// url := "https://pokeapi.co/api/v2/pokemon"
// 	// res, err := http.Get(url)
// 	// check(err)
	
// 	// defer res.Body.Close()

// 	// body, err := io.ReadAll(res.Body)
// 	// check(err)

// 	// var jsonData interface{}
// 	// err = json.Unmarshal(body, &jsonData)
// 	// check(err)

// 	// return ctx.JSON(http.StatusOK, jsonData)
// }

func check(e error) {
    if e != nil {
        panic(e)
    }
}
