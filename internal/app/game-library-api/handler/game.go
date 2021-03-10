package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Game has handler methods for dealing with games
type Game struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List returns all games
func (g *Game) List(c *gin.Context) error {
	list, err := game.List(c.Request.Context(), g.DB)

	if err != nil {
		return err
	}

	return web.Respond(c, list, http.StatusOK)
}

// Retrieve returns a game
func (g *Game) Retrieve(c *gin.Context) error {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil || id <= 0 {
		return web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}

	game, err := game.Retrieve(c.Request.Context(), g.DB, id)

	if err != nil {
		return errors.Wrapf(err, "retrieving game with id %q", idparam)
	}

	if game == nil {
		return web.NewRequestError(errors.New("Data not found"), http.StatusNotFound)
	}

	return web.Respond(c, game, http.StatusOK)
}

// Create decodes JSON and creates a new Game
func (g *Game) Create(c *gin.Context) error {
	var ng game.NewGame
	err := web.Decode(c, &ng)
	if err != nil {
		return err
	}

	game, err := game.Create(c.Request.Context(), g.DB, ng)
	if err != nil {
		return err
	}

	return web.Respond(c, game, http.StatusCreated)
}

// AddSale creates new sale for specified game
func (g *Game) AddSale(c *gin.Context) error {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 64)
	if err != nil || id <= 0 {
		return web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}

	var ns game.NewSale
	err = web.Decode(c, &ns)
	if err != nil {
		return errors.Wrap(err, "decoding new sale")
	}
	sale, err := game.AddSale(c.Request.Context(), g.DB, ns, id)

	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(c, sale, http.StatusCreated)
}

// ListSales returns sales for specified game
func (g *Game) ListSales(c *gin.Context) error {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 32)
	if err != nil || id <= 0 {
		return web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}

	list, err := game.ListSales(c.Request.Context(), g.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(c, list, http.StatusOK)
}
