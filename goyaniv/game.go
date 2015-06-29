package goyaniv

import (
	"strings"
)

type Game struct {
	Name       string    `json:"name"`
	Players    []*Player `json:"players"`
	MiddleDeck *Deck     `json:"middledeck"`
	PlayDeck   *Deck     `json:"playdeck"`
	TrashDeck  *Deck     `json:"playdeck"`
	Round      int       `json:"round"`
	Url        string    `json:"url"`
	Launched   bool      `json:"launched"`
	Yaniver    *Player   `json:"yaniver"`
	Asafed     *Player   `json:"asafed"`
	Turn       int       `json:"turn"`
	YanivAt    int       `json:"yanivat"`
}

func NewGame(gameUrl string) *Game {
	middle := NewCompleteDeck()
	middle.Shuffle()
	play := &Deck{}
	play.Add(middle.TakeCard())
	trash := &Deck{}
	return &Game{
		Name:       GetGameNameWithUrl(gameUrl),
		Players:    make([]*Player, 0),
		MiddleDeck: middle,
		PlayDeck:   play,
		TrashDeck:  trash,
		Url:        gameUrl,
		Turn:       100,
		Yaniver:    nil,
		Asafed:     nil,
		YanivAt:    5,
	}
}

func (g *Game) NewTurn() {
	g.MiddleDeck = NewCompleteDeck()
	g.MiddleDeck.Shuffle()
	g.PlayDeck = &Deck{}
	g.PlayDeck.Add(g.MiddleDeck.TakeCard())
	g.TrashDeck = &Deck{}
	g.Turn++
	g.Yaniver = nil
	g.Asafed = nil
}

func GetGameNameWithUrl(url string) string {
	return strings.Split(url, "/")[2]
}

func (g *Game) NextPlayer() {
	g.Turn++
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.Players[g.Turn%len(g.Players)]
}

func (g *Game) GetFastPlayer() *Player {
	// if first turn, no fast player
	if g.Turn == 100 {
		return g.GetCurrentPlayer()
	}
	return g.Players[(g.Turn-1)%len(g.Players)]
}

func (g *Game) GetPlayer(id string, key string) *Player {
	for _, player := range g.Players {
		if player.Id == id && player.Key == key {
			return player
		}
	}
	return nil
}

func (g *Game) AddPlayer(p *Player) {
	(*g).Players = append((*g).Players, p)
}
