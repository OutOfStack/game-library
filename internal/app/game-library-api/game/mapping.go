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
		ID:              s.ID,
		Name:            s.Name,
		GameID:          s.GameID,
		BeginDate:       s.BeginDate.String(),
		EndDate:         s.EndDate.String(),
		DiscountPercent: s.DiscountPercent,
	}
}

func (ns *NewSale) mapToGetSale(id, gameID int64) *GetSale {
	return &GetSale{
		ID:              id,
		Name:            ns.Name,
		GameID:          gameID,
		BeginDate:       ns.BeginDate,
		EndDate:         ns.EndDate,
		DiscountPercent: ns.DiscountPercent,
	}
}
