package handler

import "github.com/OutOfStack/game-library/internal/app/game-library-api/repo"

// GameResp represents game response
type GameResp struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Developer   string   `json:"developer"`
	Publisher   string   `json:"publisher"`
	ReleaseDate string   `json:"releaseDate"`
	Price       float32  `json:"price"`
	Genre       []string `json:"genre"`
	LogoURL     string   `json:"logoUrl,omitempty"`
}

// GameInfoResp represents extended game info response
type GameInfoResp struct {
	GameResp
	CurrentPrice float32 `json:"currentPrice"`
	Rating       float32 `json:"rating"`
}

// CreateGameReq represents game data we receive from user
type CreateGameReq struct {
	Name        string   `json:"name" validate:"required"`
	Developer   string   `json:"developer" validate:"required"`
	Publisher   string   `json:"-"`
	ReleaseDate string   `json:"releaseDate" validate:"date"`
	Price       float32  `json:"price" validate:"gte=0,lt=10000"`
	Genre       []string `json:"genre"`
	LogoURL     string   `json:"logoUrl"`
}

// UpdateGameReq represents model for updating information about game.\
// All fields are optional
type UpdateGameReq struct {
	Name        *string   `json:"name"`
	Developer   *string   `json:"developer" validate:"omitempty"`
	Publisher   *string   `json:"publisher" validate:"omitempty"`
	ReleaseDate *string   `json:"releaseDate" validate:"omitempty,date"`
	Price       *float32  `json:"price" validate:"omitempty,gte=0,lt=10000"`
	Genre       *[]string `json:"genre" validate:"omitempty"`
	LogoURL     *string   `json:"logoUrl"`
}

// CreateGameSaleReq represents data about game being on sale
type CreateGameSaleReq struct {
	SaleID          int64 `json:"saleId"`
	DiscountPercent uint8 `json:"discountPercent" validate:"gt=0,lte=100"`
}

// GameSaleResp represents game sale response
type GameSaleResp struct {
	GameID          int64  `json:"gameId"`
	SaleID          int64  `json:"saleId"`
	Sale            string `json:"sale"`
	DiscountPercent uint8  `json:"discountPercent"`
	BeginDate       string `json:"beginDate"`
	EndDate         string `json:"endDate"`
}

// CreateSaleReq represents sale data we receive from user
type CreateSaleReq struct {
	Name      string `json:"name"`
	BeginDate string `json:"beginDate" validate:"required,date"`
	EndDate   string `json:"endDate" validate:"required,date"`
}

// SaleResp represents sale response
type SaleResp struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

// CreateRatingReq represents rating data we receive from user
type CreateRatingReq struct {
	Rating uint8 `json:"rating" validate:"gte=1,lte=4"`
}

// RatingResp represents response to rating request
type RatingResp struct {
	GameID int64 `json:"gameId"`
	Rating uint8 `json:"rating"`
}

// UserRatingsReq represents get user ratings request
type UserRatingsReq struct {
	GameIDs []int64 `json:"gameIds"`
}

// MapToCreateRating maps CreateRatingReq to CreateRating
func mapToCreateRating(crr *CreateRatingReq, gameID int64, userID string) repo.CreateRating {
	return repo.CreateRating{
		Rating: crr.Rating,
		UserID: userID,
		GameID: gameID,
	}
}

// MapToRatingResp maps CreateRating to RatingResp
func mapToRatingResp(cr repo.CreateRating) *RatingResp {
	return &RatingResp{
		GameID: cr.GameID,
		Rating: cr.Rating,
	}
}

func mapToCreateSale(csr *CreateSaleReq) repo.CreateSale {
	return repo.CreateSale{
		Name:      csr.Name,
		BeginDate: csr.BeginDate,
		EndDate:   csr.EndDate,
	}
}

func mapCreateSaleToSaleResp(cs repo.CreateSale, id int64) *SaleResp {
	return &SaleResp{
		ID:        id,
		Name:      cs.Name,
		BeginDate: cs.BeginDate,
		EndDate:   cs.EndDate,
	}
}

func mapSaleToSaleResp(s *repo.Sale) *SaleResp {
	return &SaleResp{
		ID:        s.ID,
		Name:      s.Name,
		BeginDate: s.BeginDate.String(),
		EndDate:   s.EndDate.String(),
	}
}

func mapToCreateGameSale(cgsr *CreateGameSaleReq, gameID int64) repo.CreateGameSale {
	return repo.CreateGameSale{
		SaleID:          cgsr.SaleID,
		GameID:          gameID,
		DiscountPercent: cgsr.DiscountPercent,
	}
}

func mapToGameSaleResp(s *repo.Sale, cgs repo.CreateGameSale) *GameSaleResp {
	return &GameSaleResp{
		GameID:          cgs.GameID,
		SaleID:          cgs.SaleID,
		Sale:            s.Name,
		BeginDate:       s.BeginDate.String(),
		EndDate:         s.EndDate.String(),
		DiscountPercent: cgs.DiscountPercent,
	}
}

func mapGameSaleToGameSaleResp(gs *repo.GameSale) *GameSaleResp {
	return &GameSaleResp{
		GameID:          gs.GameID,
		SaleID:          gs.SaleID,
		Sale:            gs.Sale,
		BeginDate:       gs.BeginDate,
		EndDate:         gs.EndDate,
		DiscountPercent: gs.DiscountPercent,
	}
}

func mapToCreateGame(cgr *CreateGameReq) repo.CreateGame {
	return repo.CreateGame{
		Name:        cgr.Name,
		Developer:   cgr.Developer,
		Publisher:   cgr.Publisher,
		ReleaseDate: cgr.ReleaseDate,
		Price:       cgr.Price,
		Genre:       cgr.Genre,
		LogoURL:     cgr.LogoURL,
	}
}

func mapToUpdateGame(g *repo.Game, ugr *UpdateGameReq) repo.UpdateGame {
	var logoURL string
	if g.LogoURL.Valid {
		logoURL = g.LogoURL.String
	}
	update := repo.UpdateGame{
		ID:          g.ID,
		Name:        g.Name,
		Developer:   g.Developer,
		Publisher:   g.Publisher,
		ReleaseDate: g.ReleaseDate.String(),
		Price:       g.Price,
		LogoURL:     logoURL,
		Genre:       g.Genre,
	}

	if ugr.Name != nil {
		update.Name = *ugr.Name
	}
	if ugr.Developer != nil {
		update.Developer = *ugr.Developer
	}
	if ugr.Publisher != nil {
		update.Publisher = *ugr.Publisher
	}
	if ugr.ReleaseDate != nil {
		update.ReleaseDate = *ugr.ReleaseDate
	}
	if ugr.Price != nil {
		update.Price = *ugr.Price
	}
	if ugr.Genre != nil {
		update.Genre = *ugr.Genre
	}
	if ugr.LogoURL != nil && *ugr.LogoURL != "" {
		update.LogoURL = *ugr.LogoURL
	}

	return update
}

func mapToGameInfoResp(g *repo.GameExt) *GameInfoResp {
	return &GameInfoResp{
		GameResp: GameResp{
			ID:          g.ID,
			Name:        g.Name,
			Developer:   g.Developer,
			Publisher:   g.Publisher,
			ReleaseDate: g.ReleaseDate.String(),
			Price:       g.Price,
			Genre:       []string(g.Genre),
			LogoURL:     g.LogoURL.String,
		},
		CurrentPrice: g.CurrentPrice,
		Rating:       g.Rating,
	}
}

func mapToGameResp(cg repo.CreateGame, id int64) *GameResp {
	return &GameResp{
		ID:          id,
		Name:        cg.Name,
		Developer:   cg.Developer,
		Publisher:   cg.Publisher,
		ReleaseDate: cg.ReleaseDate,
		Price:       cg.Price,
		Genre:       cg.Genre,
		LogoURL:     cg.LogoURL,
	}
}

func mapUpdateGameToGameResp(ug repo.UpdateGame) *GameResp {
	return &GameResp{
		ID:          ug.ID,
		Name:        ug.Name,
		Developer:   ug.Developer,
		Publisher:   ug.Publisher,
		ReleaseDate: ug.ReleaseDate,
		Price:       ug.Price,
		Genre:       ug.Genre,
		LogoURL:     ug.LogoURL,
	}
}
