package game

func (g *Game) mapToGetGame() *GetGame {
	return &GetGame{
		ID:          g.ID,
		Name:        g.Name,
		Developer:   g.Developer,
		ReleaseDate: g.ReleaseDate.String(),
		Price:       g.Price,
		Genre:       []string(g.Genre),
	}
}

func (ng *NewGame) mapToGetGame(id int64) *GetGame {
	return &GetGame{
		ID:          id,
		Name:        ng.Name,
		Developer:   ng.Developer,
		ReleaseDate: ng.ReleaseDate,
		Price:       ng.Price,
		Genre:       ng.Genre,
	}
}

func (s *Sale) mapToGetSale() *GetSale {
	return &GetSale{
		ID:        s.ID,
		Name:      s.Name,
		BeginDate: s.BeginDate.String(),
		EndDate:   s.EndDate.String(),
	}
}

func (ns *NewSale) mapToGetSale(id int64) *GetSale {
	return &GetSale{
		ID:        id,
		Name:      ns.Name,
		BeginDate: ns.BeginDate,
		EndDate:   ns.EndDate,
	}
}

func (gs *GameSale) mapToGetGameSale() *GetGameSale {
	return &GetGameSale{
		GameID:          gs.GameID,
		SaleID:          gs.SaleID,
		Sale:            gs.Sale,
		BeginDate:       gs.BeginDate,
		EndDate:         gs.EndDate,
		DiscountPercent: gs.DiscountPercent,
	}
}

func (ngs *NewGameSale) mapToGetGameSale(sale *GetSale, gameId int64) *GetGameSale {
	return &GetGameSale{
		GameID:          gameId,
		SaleID:          ngs.SaleID,
		Sale:            sale.Name,
		BeginDate:       sale.BeginDate,
		EndDate:         sale.EndDate,
		DiscountPercent: ngs.DiscountPercent,
	}
}
