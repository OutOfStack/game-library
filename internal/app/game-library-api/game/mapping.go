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
		LogoUrl:     g.LogoUrl.String,
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
			LogoUrl:     g.LogoUrl.String,
		},
		CurrentPrice: g.CurrentPrice,
		Rating:       g.Rating,
	}
}

// MapToGameResp maps CreateGameReq to GameResp
func (ng *CreateGameReq) MapToGameResp(id int64) *GameResp {
	return &GameResp{
		ID:          id,
		Name:        ng.Name,
		Developer:   ng.Developer,
		Publisher:   ng.Publisher,
		ReleaseDate: ng.ReleaseDate,
		Price:       ng.Price,
		Genre:       ng.Genre,
		LogoUrl:     ng.LogoUrl,
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

// MapToSaleResp maps CreateSaleReq to SaleResp
func (ns *CreateSaleReq) MapToSaleResp(id int64) *SaleResp {
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

// MapToGameSale maps CreateGameSaleReq and Sale to GameSale
func (ngs *CreateGameSaleReq) MapToGameSale(sale *Sale, gameID int64) *GameSale {
	return &GameSale{
		GameID:          gameID,
		SaleID:          ngs.SaleID,
		Sale:            sale.Name,
		BeginDate:       sale.BeginDate.String(),
		EndDate:         sale.EndDate.String(),
		DiscountPercent: ngs.DiscountPercent,
	}
}
