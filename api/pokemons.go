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
	var args models.ListPokemonsOffsetParams

	limitQuery := c.QueryParam("limit")
	limit, err := strconv.ParseInt(limitQuery, 10, 32)
	if err != nil {
		limit = 10
	}
	offsetQuery := c.QueryParam("offset")
	offset, err := strconv.ParseInt(offsetQuery, 10, 64)
	if err != nil {
		offset = 0
	}

	args.Limit = limit
	args.Offset = offset

	// pokemons, err := models.QUERIES.ListPokemons(models.CTX)
	pokemons, err := models.QUERIES.ListPokemonsOffset(models.CTX, args)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	data := DataList{
		Count: len(pokemons),
		Limit: int(limit),
		Offset: int(offset),
		// TODO fill Previous & Next
		// Previous: `/pokemons?limit=3&offset=3`,
		// Next: `/pokemons?limit=5&offset=3`,
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
