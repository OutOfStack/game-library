package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/game"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/web"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Game has handler methods for dealing with games
type Game struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List returns all games
func (g *Game) List(c *gin.Context) error {
	list, err := game.List(g.DB)

	if err != nil {
		return err
	}

	return web.Respond(c, list, http.StatusOK)
}

// Retrieve returns a game
func (g *Game) Retrieve(c *gin.Context) error {
	idparam := c.Param("id")
	id, err := strconv.ParseUint(idparam, 10, 64)
	if err != nil {
		return err
	}

	game, err := game.Retrieve(g.DB, id)

	if err != nil {
		return err
	}

	if game == nil {
		return web.NewErrorRequest(errors.New("Data not found"), http.StatusNotFound)
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

	game, err := game.Create(g.DB, pm)
	if err != nil {
		return err
	}

	return web.Respond(c, game, http.StatusCreated)
}
