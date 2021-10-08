package game

// MapToGetGame maps Game to GetGame
func (g *Game) MapToGetGame() *GetGame {
	return &GetGame{
		ID:          g.ID,
		Name:        g.Name,
		Developer:   g.Developer,
		Publisher:   g.Publisher,
		ReleaseDate: g.ReleaseDate.String(),
		Price:       g.Price,
		Genre:       []string(g.Genre),
	}
}

// MapToGetGameInfo maps GameInfo to GetGameInfo
func (g *GameInfo) MapToGetGameInfo() *GetGameInfo {
	return &GetGameInfo{
		GetGame: GetGame{
			ID:          g.ID,
			Name:        g.Name,
			Developer:   g.Developer,
			Publisher:   g.Publisher,
			ReleaseDate: g.ReleaseDate.String(),
			Price:       g.Price,
			Genre:       []string(g.Genre),
		},
		CurrentPrice: g.CurrentPrice,
		Rating:       g.Rating,
	}
}

// MapToGetGame maps NewGame to GetGame
func (ng *CreateGame) MapToGetGame(id int64) *GetGame {
	return &GetGame{
		ID:          id,
		Name:        ng.Name,
		Developer:   ng.Developer,
		Publisher:   ng.Publisher,
		ReleaseDate: ng.ReleaseDate,
		Price:       ng.Price,
		Genre:       ng.Genre,
	}
}

// MapToGetSale maps Sale to GetSale
func (s *Sale) MapToGetSale() *GetSale {
	return &GetSale{
		ID:        s.ID,
		Name:      s.Name,
		BeginDate: s.BeginDate.String(),
		EndDate:   s.EndDate.String(),
	}
}

// MapToGetSale maps NewSale to GetSale
func (ns *CreateSale) MapToGetSale(id int64) *GetSale {
	return &GetSale{
		ID:        id,
		Name:      ns.Name,
		BeginDate: ns.BeginDate,
		EndDate:   ns.EndDate,
	}
}

// MapToGetGameSale maps GameSale to GetGameSale
func (gs *GameSale) MapToGetGameSale() *GetGameSale {
	return &GetGameSale{
		GameID:          gs.GameID,
		SaleID:          gs.SaleID,
		Sale:            gs.Sale,
		BeginDate:       gs.BeginDate,
		EndDate:         gs.EndDate,
		DiscountPercent: gs.DiscountPercent,
	}
}

// MapToGameSale maps CreateGameSale and Sale to GameSale
func (ngs *CreateGameSale) MapToGameSale(sale *Sale, gameID int64) *GameSale {
	return &GameSale{
		GameID:          gameID,
		SaleID:          ngs.SaleID,
		Sale:            sale.Name,
		BeginDate:       sale.BeginDate.String(),
		EndDate:         sale.EndDate.String(),
		DiscountPercent: ngs.DiscountPercent,
	}
}
