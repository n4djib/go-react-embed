package api

import (
	"go-react-embed/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RegisterPokemonsHandlers(e *echo.Group) {
	e.GET("/pokemons", getPokemonsHandler)
	e.GET("/pokemons/:id", getPokemonHandler)
	e.GET("/pokemons/name/:name", getPokemonByNameHandler)
}

func getPokemonsHandler(c echo.Context) error {
	pokemons, err := models.QUERIES.ListPokemons(models.CTX)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := DataList{
		Count: len(pokemons),
		Data: pokemons,
	}
	return c.JSON(http.StatusOK, data)
}

func getPokemonHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{
			Error: err.Error(),
		})
	}
	pokemon, err := models.QUERIES.GetPokemon(models.CTX, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := Data{
		Data: pokemon,
	}
	return c.JSON(http.StatusOK, data)
}

func getPokemonByNameHandler(c echo.Context) error {
	name := c.Param("name")
	pokemon, err := models.QUERIES.GetPokemonByName(models.CTX, name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := Data{
		Data: pokemon,
	}
	return c.JSON(http.StatusOK, data)
}
