package api

import (
	"context"
	"database/sql"
	"go-react-embed/models"
	"go-react-embed/rbac"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PokemonHandler struct {
	CTX     context.Context
	QUERIES *models.Queries
	DB      *sql.DB
	RBAC    rbac.RBAC
}

func (h PokemonHandler) RegisterHandlers(e *echo.Group) {
	e.GET("/pokemons", h.getPokemonsHandler)
	e.GET("/pokemons/:id", h.getPokemonHandler)
	e.GET("/pokemons/name/:name", h.getPokemonByNameHandler)
}

// @Summary Get All Pokemons
// @Description List all pokemons (limit & offset)
// @Tags Pokemons
// @Param limit query int false "Limit: default 10"
// @Param offset query int false "Offset: default 0"
// @Produce json
// @Success 200 {string} string "ok"
// @Router /api/pokemons [get]
func (h PokemonHandler) getPokemonsHandler(c echo.Context) error {
	// get current user and roles from context

	// check if user has permission
	// IsAllowed(user, "permission") bool

	// TODO read this and do some changes to reduce the code
	// https://dev.to/geoff89/deploying-a-golang-restful-api-with-gin-sqlc-and-postgresql-1lbl#implementing-crud-in-golang-rest-api

	
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

	var args models.ListPokemonsOffsetParams
	args.Limit = limit
	args.Offset = offset
	// args := &models.ListPokemonsOffsetParams{
    //     Limit:  limit,
    //     Offset: offset,
    // }

	pokemonsCount, err := h.QUERIES.GetPokemonsCount(h.CTX)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	pokemons, err := h.QUERIES.ListPokemonsOffset(h.CTX, args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, echo.Map{
		"count": pokemonsCount[0],
		"limit": int(limit),
		"offset": int(offset),
		"result": pokemons,
	})
}


// TODO i should remove this
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
func (h PokemonHandler) getPokemonHandler(c echo.Context) error {
	// ctxUser := GetUserFromContext(c)
	// if ctxUser.ID == 0 {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, "Not Authenticated")
	// }
	
	// user := rbac.Map{
	// 	"id": ctxUser.ID, 
	// 	"name": ctxUser.Name, 
	// 	"roles": ctxUser.Roles,
	// }
	// resource := rbac.Map{"id": 3, "title": "tutorial", "owner": 3, "list": []int{1, 2, 3, 4, 5, 6}}

	// allowed, err := h.RBAC.IsAllowed(user, resource, "edit_user")
	// if err != nil {
	// 	fmt.Println("++++ error: ", err.Error())
	// }
	// fmt.Println("-allowed to get pokemon:", allowed)

	param := c.Param("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "id param must be integer")
	}
	pokemon, err := h.QUERIES.GetPokemon(h.CTX, id)
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
func (h PokemonHandler) getPokemonByNameHandler(c echo.Context) error {
	name := c.Param("name")
	pokemon, err := h.QUERIES.GetPokemonByName(h.CTX, name)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "no rows in result set")
	}
	return c.JSON(http.StatusOK, PokemonResult{
		Result: pokemon,
	})
}
