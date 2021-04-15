package handler

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/OutOfStack/game-library/internal/app/game-library-api/game"
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
	list, err := repo.List(c.Request.Context(), g.DB)

	if err != nil {
		return errors.Wrap(err, "getting games list")
	}

	return web.Respond(c, list, http.StatusOK)
}

// Retrieve returns a game
func (g *Game) Retrieve(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	game, err := repo.Retrieve(c.Request.Context(), g.DB, id)

	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "retrieving game with id %q", id)
	}

	return web.Respond(c, game, http.StatusOK)
}

// Create decodes JSON and creates a new Game
func (g *Game) Create(c *gin.Context) error {
	var ng repo.NewGame
	err := web.Decode(c, &ng)
	if err != nil {
		return errors.Wrap(err, "decoding new game")
	}

	game, err := repo.Create(c.Request.Context(), g.DB, ng)
	if err != nil {
		return errors.Wrap(err, "adding new game")
	}

	return web.Respond(c, game, http.StatusCreated)
}

// AddSale creates new sale for specified game
func (g *Game) AddSale(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	var ns repo.NewSale
	err = web.Decode(c, &ns)
	if err != nil {
		return errors.Wrap(err, "decoding new sale")
	}
	sale, err := repo.AddSale(c.Request.Context(), g.DB, ns, id)

	if err != nil {
		return errors.Wrap(err, "adding new sale")
	}

	return web.Respond(c, sale, http.StatusCreated)
}

// ListSales returns sales for specified game
func (g *Game) ListSales(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	list, err := repo.ListSales(c.Request.Context(), g.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(c, list, http.StatusOK)
}

// Update updates specified game
func (g *Game) Update(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}
	var update repo.UpdateGame
	if err := web.Decode(c, &update); err != nil {
		return errors.Wrap(err, "decoding game update")
	}
	err = repo.Update(c.Request.Context(), g.DB, id, update)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "updating game with id %q", id)
	}

	return web.Respond(c, nil, http.StatusNoContent)
}

// Delete removes specified game
func (g *Game) Delete(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	err = repo.Delete(c.Request.Context(), g.DB, id)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "deleting game with id %q", id)
	}

	return web.Respond(c, nil, http.StatusNoContent)
}

func getIdParam(c *gin.Context) (int64, error) {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 32)
	if err != nil || id <= 0 {
		return 0, web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}
	return id, err
}
