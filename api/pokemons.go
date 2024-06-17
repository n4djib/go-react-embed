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

// @Summary Get All Pokemons
// @Description List all pokemons (limit & offset)
// @Tags Pokemons
// @Param limit query int false "Limit: default 10"
// @Param offset query int false "Offset: default 0"
// @Produce json
// @Success 200 {string} string "ok"
// @Router /api/pokemons [get]
func getPokemonsHandler(c echo.Context) error {
	// get current user and roles from context

	// check if user has permission
	// IsAllowed(user, "permission") bool

	
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

	allPokemonsNames, err := models.QUERIES.ListPokemonsNames(models.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	pokemons, err := models.QUERIES.ListPokemonsOffset(models.CTX, args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"count": len(allPokemonsNames),
		"limit": int(limit),
		"offset": int(offset),
		"result": pokemons,
	})
}



type PokemonResult struct {
	Result models.Pokemon `json:"result"`
}

// @Summary Get Pokemon by ID
// @Description get pokemon by ID as param path
// @Tags Pokemons
// @Param id path int true "ID of a pokemon"
// @Produce json
// @Success 200 {object} PokemonResult
// @Router /api/pokemons/{id} [get]
func getPokemonHandler(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id param must be integer")
	}
	pokemon, err := models.QUERIES.GetPokemon(models.CTX, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, PokemonResult{
		Result: pokemon,
	})
}

// @Summary Get Pokemon by Name
// @Description get pokemon by Name as param path
// @Tags Pokemons
// @Param name path string true "Name of a pokemon"
// @Produce json
// @Success 200 {object} PokemonResult
// @Router /api/pokemons/name/{name} [get]
func getPokemonByNameHandler(c echo.Context) error {
	name := c.Param("name")
	pokemon, err := models.QUERIES.GetPokemonByName(models.CTX, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, PokemonResult{
		Result: pokemon,
	})
}
