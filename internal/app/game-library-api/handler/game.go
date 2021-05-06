package handler

import (
	"context"
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

// List godoc
// @Summary List all games
// @Description returns all games
// @ID get-all-games
// @Produce json
// @Success 200 {array} game.GetGame
// @Failure 500 {object} web.ErrorResponse
// @Router /games [get]
func (g *Game) List(ctx context.Context, c *gin.Context) error {
	list, err := repo.List(c.Request.Context(), g.DB)

	if err != nil {
		return errors.Wrap(err, "getting games list")
	}

	return web.Respond(ctx, c, list, http.StatusOK)
}

// Retrieve godoc
// @Summary Show a game
// @Description returns game by ID
// @ID get-game-by-id
// @Produce json
// @Param id path int64 true "Game ID"
// @Success 200 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (g *Game) Retrieve(ctx context.Context, c *gin.Context) error {
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

	return web.Respond(ctx, c, game, http.StatusOK)
}

// Create godoc
// @Summary Create a game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce  json
// @Param  game body game.NewGame true "create game"
// @Success 200 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (g *Game) Create(ctx context.Context, c *gin.Context) error {
	var ng repo.NewGame
	err := web.Decode(c, &ng)
	if err != nil {
		return errors.Wrap(err, "decoding new game")
	}

	game, err := repo.Create(c.Request.Context(), g.DB, ng)
	if err != nil {
		return errors.Wrap(err, "adding new game")
	}

	return web.Respond(ctx, c, game, http.StatusCreated)
}

// Update godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game-by-id
// @Accept  json
// @Produce  json
// @Param  id path int64 true "Game ID"
// @Param  game body game.UpdateGame true "update game"
// @Success 200 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (g *Game) Update(ctx context.Context, c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}
	var update repo.UpdateGame
	if err := web.Decode(c, &update); err != nil {
		return errors.Wrap(err, "decoding game update")
	}
	game, err := repo.Update(c.Request.Context(), g.DB, id, update)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "updating game with id %q", id)
	}

	return web.Respond(ctx, c, game, http.StatusOK)
}

// Delete godoc
// @Summary Delete a game
// @Description deletes game by ID
// @ID delete-game-by-id
// @Accept  json
// @Produce  json
// @Param  id path int64 true "Game ID"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (g *Game) Delete(ctx context.Context, c *gin.Context) error {
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

	return web.Respond(ctx, c, nil, http.StatusNoContent)
}

// AddSale godoc
// @Summary Create a sale
// @Description Creates new sale
// @ID create-sale
// @Accept  json
// @Produce  json
// @Param id path int64 true "Game ID"
// @Param  sale body game.NewSale true "create sale"
// @Success 200 {object} game.GetSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [post]
func (g *Game) AddSale(ctx context.Context, c *gin.Context) error {
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
		if errors.Is(err, repo.ErrNotFound) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "adding new sale for gameID %q", id)
	}

	return web.Respond(ctx, c, sale, http.StatusCreated)
}

// ListSales godoc
// @Summary List all sales
// @Description Returns all sales
// @ID get-sales-for-game
// @Produce json
// @Param id path int64 true "Game ID"
// @Success 200 {array} game.GetSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [get]
func (g *Game) ListSales(ctx context.Context, c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	list, err := repo.ListSales(c.Request.Context(), g.DB, id)
	if err != nil {
		return errors.Wrap(err, "getting sales list")
	}

	return web.Respond(ctx, c, list, http.StatusOK)
}

func getIdParam(c *gin.Context) (int64, error) {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 32)
	if err != nil || id <= 0 {
		return 0, web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}
	return id, err
}
