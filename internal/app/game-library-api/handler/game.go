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

// List godoc
// @Summary List all games
// @Description returns all games
// @ID get-all-games
// @Produce json
// @Success 200 {array} game.GetGame
// @Failure 500 {object} web.ErrorResponse
// @Router /games [get]
func (g *Game) List(c *gin.Context) error {
	list, err := repo.List(c, g.DB)

	if err != nil {
		return errors.Wrap(err, "getting games list")
	}

	getGames := []repo.GetGame{}
	for _, g := range list {
		getGames = append(getGames, *g.MapToGetGame())
	}

	return web.Respond(c, getGames, http.StatusOK)
}

// Retrieve godoc
// @Summary Show a game
// @Description returns game by ID
// @ID get-game-by-id
// @Produce json
// @Param 	id  path int64 true "Game ID"
// @Success 200 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [get]
func (g *Game) Retrieve(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	game, err := repo.Retrieve(c, g.DB, id)

	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "retrieving game with id %v", id)
	}

	getGame := game.MapToGetGame()

	return web.Respond(c, getGame, http.StatusOK)
}

// Create godoc
// @Summary Create a game
// @Description creates new game
// @ID create-game
// @Accept  json
// @Produce json
// @Param   game body game.CreateGame true "create game"
// @Success 201 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games [post]
func (g *Game) Create(c *gin.Context) error {
	var cg repo.CreateGame
	err := web.Decode(c, &cg)
	if err != nil {
		return errors.Wrap(err, "decoding new game")
	}

	gameId, err := repo.Create(c, g.DB, cg)
	if err != nil {
		return errors.Wrap(err, "adding new game")
	}

	getGame := cg.MapToGetGame(gameId)

	return web.Respond(c, getGame, http.StatusCreated)
}

// Update godoc
// @Summary Update game
// @Description updates game by ID
// @ID update-game-by-id
// @Accept  json
// @Produce json
// @Param  	id   path int64 			true "Game ID"
// @Param  	game body game.UpdateGame 	true "update game"
// @Success 200 {object} game.GetGame
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [patch]
func (g *Game) Update(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}
	var update repo.UpdateGame
	if err := web.Decode(c, &update); err != nil {
		return errors.Wrap(err, "decoding game update")
	}
	game, err := repo.Update(c, g.DB, id, update)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "updating game with id %v", id)
	}

	getGame := game.MapToGetGame()

	return web.Respond(c, getGame, http.StatusOK)
}

// Delete godoc
// @Summary Delete a game
// @Description deletes game by ID
// @ID delete-game-by-id
// @Accept  json
// @Produce json
// @Param  	id path int64 true "Game ID"
// @Success 204
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id} [delete]
func (g *Game) Delete(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}

	err = repo.Delete(c, g.DB, id)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "deleting game with id %v", id)
	}

	return web.Respond(c, nil, http.StatusNoContent)
}

// AddGameOnSale godoc
// @Summary Add game on sale
// @Description adds game on sale
// @ID add-game-on-sale
// @Accept  json
// @Produce json
// @Param  id 		path int64 				 true "Game ID"
// @Param  gamesale body game.CreateGameSale true "game sale"
// @Success 200 {object} game.GetGameSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [post]
func (g *Game) AddGameOnSale(c *gin.Context) error {
	id, err := getIdParam(c)
	if err != nil {
		return err
	}
	var cgs repo.CreateGameSale
	if err := web.Decode(c, &cgs); err != nil {
		return errors.Wrap(err, "decoding game sale")
	}
	gameSale, err := repo.AddGameOnSale(c, g.DB, id, cgs)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrap(err, "add game on sale")
	}

	getGameSale := gameSale.MapToGetGameSale()

	return web.Respond(c, getGameSale, http.StatusOK)
}

// ListGameSales godoc
// @Summary List game sales
// @Description returns sales for specified game
// @ID get-game-sales-by-game-id
// @Produce json
// @Param 	id  path int64 true "Game ID"
// @Success 200 {array}  game.GetGameSale
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /games/{id}/sales [get]
func (g *Game) ListGameSales(c *gin.Context) error {
	gameId, err := getIdParam(c)
	if err != nil {
		return err
	}

	gameSales, err := repo.ListGameSales(c, g.DB, gameId)
	if err != nil {
		if errors.As(err, &repo.ErrNotFound{}) {
			return web.NewRequestError(err, http.StatusNotFound)
		}
		return errors.Wrapf(err, "retrieving sales for game")
	}

	getGameSales := []repo.GetGameSale{}
	for _, gs := range gameSales {
		getGameSales = append(getGameSales, *gs.MapToGetGameSale())
	}

	return web.Respond(c, getGameSales, http.StatusOK)
}

func getIdParam(c *gin.Context) (int64, error) {
	idparam := c.Param("id")
	id, err := strconv.ParseInt(idparam, 10, 32)
	if err != nil || id <= 0 {
		return 0, web.NewRequestError(errors.New("Invalid id"), http.StatusBadRequest)
	}
	return id, err
}
