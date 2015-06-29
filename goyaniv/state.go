package goyaniv

import (
	"encoding/json"
)

type LastAction struct {
	Player   string `json:"player"`
	TakeCard int    `json:"takecard"` //-1 = Yaniv, -2 = Asaf
}

type StatePlayer struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Me         bool   `json:"me"`
	Playing    bool   `json:"playing"`
	Connected  bool   `json:"connected"`
	Yaniver    bool   `json:"yaniver"`
	Ready      bool   `json:"ready"`
	Lost       bool   `json:"lost"`
	Score      int    `json:"score"`
	DeckWeight int    `json:"deckweight"`
	Deck       Deck   `json:"deck"`
}

type State struct {
	PlayDeck   Deck       `json:"playdeck"`
	LastAction LastAction `json:"lastaction"`

	Round      int           `json:"round"`
	Started    bool          `json:"started"`
	Terminated bool          `json:"terminated"`
	Players    []StatePlayer `json:"players"`
	Error      string        `json:"error"`
}

func NewStatePlayer(p *Player) StatePlayer {
	return StatePlayer{
		Name:       p.Name,
		Id:         p.Id,
		Me:         true,
		Playing:    false,
		Connected:  true,
		Yaniver:    false,
		Ready:      p.Ready,
		Lost:       false,
		Score:      p.Score,
		DeckWeight: p.Deck.Weight(),
		Deck:       *p.Deck,
	}
}

func NewStateError(g *Game, p *Player, error string) *State {
	state := NewState(g, p)
	state.Error = error
	return state
}

func (sp *StatePlayer) HideInfos() {
	HiddenCard := &Card{}
	sp.Me = false
	sp.DeckWeight = 2
	for i, _ := range sp.Deck {
		(*sp).Deck[i] = HiddenCard
	}
}
func NewState(g *Game, p *Player) *State {
	stateplayers := make([]StatePlayer, 0)
	for _, player := range g.Players {
		sp := NewStatePlayer(player)
		if sp.Id != p.Id {
			sp.HideInfos()
		}
		stateplayers = append(stateplayers, sp)
	}
	return &State{
		PlayDeck:   *g.PlayDeck,
		LastAction: LastAction{},
		Round:      g.Round,
		Players:    stateplayers,
	}
}

func (s State) Json() (string, []byte) {
	data, _ := json.Marshal(s)
	return string(data), data
}
