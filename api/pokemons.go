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

	// TODO return errors[] if params are bad
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
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	pokemons, err := models.QUERIES.ListPokemonsOffset(models.CTX, args)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{
			Error: err.Error(),
		})
	}
	// data := DataList{
	// 	Count:  len(allPokemonsNames),
	// 	Limit:  int(limit),
	// 	Offset: int(offset),
	// 	Data:   pokemons,
	// }
	// return c.JSON(http.StatusOK, data)
	return c.JSON(http.StatusOK, echo.Map{
		"count": len(allPokemonsNames),
		"limit": int(limit),
		"offset": int(offset),
		"result": pokemons,
	})
}

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
	// result := Result{
	// 	Result: pokemon,
	// }
	// return c.JSON(http.StatusOK, result)
	return c.JSON(http.StatusOK, echo.Map{
		"result": pokemon,
	})
}

func getPokemonByNameHandler(c echo.Context) error {
	name := c.Param("name")
	pokemon, err := models.QUERIES.GetPokemonByName(models.CTX, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	// data := Data{
	// 	Data: pokemon,
	// }
	// return c.JSON(http.StatusOK, data)
	return c.JSON(http.StatusOK, echo.Map{
		"result": pokemon,
	})
}
