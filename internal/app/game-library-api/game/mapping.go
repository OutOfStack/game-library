package game

// MapToGameResp maps Game to GameResp
func (g *Game) MapToGameResp() *GameResp {
	return &GameResp{
		ID:          g.ID,
		Name:        g.Name,
		Developer:   g.Developer,
		Publisher:   g.Publisher,
		ReleaseDate: g.ReleaseDate.String(),
		Price:       g.Price,
		Genre:       []string(g.Genre),
	}
}

// MapToGameInfoResp maps GameInfo to GameInfoResp
func (g *GameInfo) MapToGameInfoResp() *GameInfoResp {
	return &GameInfoResp{
		GameResp: GameResp{
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

// MapToGameResp maps NewGame to GameResp
func (ng *CreateGame) MapToGameResp(id int64) *GameResp {
	return &GameResp{
		ID:          id,
		Name:        ng.Name,
		Developer:   ng.Developer,
		Publisher:   ng.Publisher,
		ReleaseDate: ng.ReleaseDate,
		Price:       ng.Price,
		Genre:       ng.Genre,
	}
}

// MapToSaleResp maps Sale to SaleResp
func (s *Sale) MapToSaleResp() *SaleResp {
	return &SaleResp{
		ID:        s.ID,
		Name:      s.Name,
		BeginDate: s.BeginDate.String(),
		EndDate:   s.EndDate.String(),
	}
}

// MapToSaleResp maps NewSale to SaleResp
func (ns *CreateSale) MapToSaleResp(id int64) *SaleResp {
	return &SaleResp{
		ID:        id,
		Name:      ns.Name,
		BeginDate: ns.BeginDate,
		EndDate:   ns.EndDate,
	}
}

// MapToGameSaleResp maps GameSale to GameSaleResp
func (gs *GameSale) MapToGameSaleResp() *GameSaleResp {
	return &GameSaleResp{
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
