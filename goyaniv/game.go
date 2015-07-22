package goyaniv

import (
	"fmt"
	"sort"
	"strings"
)

type Game struct {
	Name          string    `json:"name"`
	Players       []*Player `json:"players"`
	ReferenceDeck *Deck     `json:"referencedeck"`
	MiddleDeck    *Deck     `json:"middledeck"`
	PlayDeck      *Deck     `json:"playdeck"`
	TrashDeck     *Deck     `json:"playdeck"`
	Round         int       `json:"round"`
	Url           string    `json:"url"`
	Launched      bool      `json:"launched"`
	Turn          int       `json:"turn"`
	YanivAt       int       `json:"yanivat"`
	LastLog       *Log      `json:"lastlog"`
}

func (g *Game) CardFromReference(id int) Card {
	for _, card := range *g.ReferenceDeck {
		if card.Id == id {
			return *card
		}
	}
	return Card{}
}

func NewGame(gameUrl string) *Game {
	referencedeck := NewCompleteDeck()
	middle := NewCompleteDeck()
	middle.Shuffle()
	play := &Deck{}
	play.Add(middle.TakeCard())
	trash := &Deck{}
	return &Game{
		Name:          GetGameNameWithUrl(gameUrl),
		Players:       make([]*Player, 0),
		ReferenceDeck: referencedeck,
		MiddleDeck:    middle,
		PlayDeck:      play,
		TrashDeck:     trash,
		Url:           gameUrl,
		Turn:          100,
		YanivAt:       5,
		LastLog:       &Log{},
	}
}

func (g *Game) IsTurnFinished() bool {
	for _, player := range g.Players {
		if player.WantsAsaf == "noanswer" {
			return false
		}
	}
	return true
}

func (g *Game) PlayersWantsAsaf() ListPlayer {
	pwa := ListPlayer{}
	for _, player := range g.Players {
		if player.WantsAsaf == "yes" {
			pwa = append(pwa, player)
		}
	}
	return pwa
}

func (g *Game) GetReferenceCardById(id int) *Card {
	for _, card := range *g.ReferenceDeck {
		if card.Id == id {
			return card
		}
	}
	return nil
}

func (g *Game) AllPlayersAnswered() bool {
	for _, player := range g.Players {
		fmt.Println(player.WantsAsaf)
		if player.WantsAsaf == "noanswer" {
			return false
		}
	}
	return true
}

func (g *Game) UpdateScores() {
	if g.AllPlayersAnswered() {
		pwa := g.PlayersWantsAsaf()
		sort.Sort(pwa)
		for i, player := range pwa {
			if i > 0 {
				player.Score = player.Score + i*30 - player.Deck.Weight()
			} else {
				player.Score = player.Score - player.Deck.Weight()
				// the winner begins
				for i, ptmp := range g.Players {
					if ptmp == player {
						g.Turn = i + len(g.Players)
					}
				}
			}
		}

		for _, player := range g.Players {
			player.Score = player.Score + player.Deck.Weight()
			// Halved
			if player.Score%50 == 0 {
				if player.WantsAsaf == "no" {
					player.Score = player.Score - 50
				}
			}
			if player.Score >= 200 {
				player.State = "loose"
			}
		}

		g.NewTurn()
	}
}

func (g *Game) GetAsafRank() int {
	var i int
	for _, player := range g.Players {
		if player.Asaf > i {
			i = player.Asaf
		}
	}
	return i + 1
}

func (g *Game) NewTurn() {
	g.MiddleDeck = NewCompleteDeck()
	g.MiddleDeck.Shuffle()
	g.PlayDeck = &Deck{}
	g.PlayDeck.Add(g.MiddleDeck.TakeCard())
	g.TrashDeck = &Deck{}
	for _, player := range g.Players {
		player.Yaniv = false
		player.Asaf = 0
		player.WantsAsaf = "noanswer"
		player.Deck = &Deck{}
		for i := 0; i < 5; i++ {
			player.Deck.Add(g.MiddleDeck.TakeCard())
		}
	}
}

func GetGameNameWithUrl(url string) string {
	return strings.Split(url, "/")[2]
}

func (g *Game) NextPlayer() {
	if len(*g.MiddleDeck) == 0 {
		g.TrashDeck.Shuffle()
		g.MiddleDeck.AddDeck(g.TrashDeck)
		fmt.Println("Middle Deck reset")
	}
	g.Turn++
}

func (g *Game) GetCurrentPlayer() *Player {
	return g.Players[g.Turn%len(g.PlayersPlaying())]
}

func (g *Game) PlayersPlaying() []*Player {
	pl := make([]*Player, 0)
	for _, player := range g.Players {
		if player.State == "playing" {
			pl = append(pl, player)
		}
	}
	return pl
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
