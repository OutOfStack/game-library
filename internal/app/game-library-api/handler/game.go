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
	id, err := strconv.ParseUint(idparam, 10, 32)
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
	var pm game.PostModel
	err := web.Decode(c, &pm)
	if err != nil {
		return err
	}

	game, err := game.Create(c.Request.Context(), g.DB, pm)
	if err != nil {
		return err
	}

	return web.Respond(c, game, http.StatusCreated)
}
