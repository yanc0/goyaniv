package goyaniv

type Log struct {
	PlayerName string `json:"playername"`
	Action     string `json:"action"`
	TakeCard   Card   `json:"takecard"`
	PutCards   []Card `json:"putcards"`
	Option     string `json:"option"`
}

func (a *Action) ToLog(g *Game) {
	PutCards := make([]Card, 0)
	TakeCard := g.CardFromReference(a.TakeCard)
	for _, id := range a.PutCards {
		PutCards = append(PutCards, g.CardFromReference(id))
	}
	g.LastLog = &Log{a.PlayerName, a.Name, TakeCard, PutCards, a.Option}
}
