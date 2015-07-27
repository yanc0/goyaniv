package goyaniv

import (
	"encoding/json"
	"sort"
)

type StatePlayer struct {
	Name       string `json:"name"`
	Id         string `json:"id"`
	Me         bool   `json:"me"`
	Spectator  bool   `json:"spectator"`
	Playing    bool   `json:"playing"`
	Connected  bool   `json:"connected"`
	Yaniver    bool   `json:"yaniver"`
	Asafer     bool   `json:"asafer"`
	Ready      bool   `json:"ready"`
	Lost       bool   `json:"lost"`
	Score      int    `json:"score"`
	DeckWeight int    `json:"deckweight"`
	Deck       Deck   `json:"deck"`
}

type State struct {
	PlayDeck   Deck          `json:"playdeck"`
	LastLog    Log           `json:"lastlog"`
	Round      int           `json:"round"`
	Started    bool          `json:"started"`
	Terminated bool          `json:"terminated"`
	Players    []StatePlayer `json:"players"`
	Error      string        `json:"error"`
}

func NewStatePlayer(p *Player, playing bool) StatePlayer {
	sort.Sort(*p.Deck)

	return StatePlayer{
		Name:       p.Name,
		Id:         p.Id,
		Me:         true,
		Spectator:  p.State == "spectator",
		Playing:    playing,
		Connected:  p.Connected,
		Yaniver:    p.Yaniv,
		Asafer:     (p.Asaf > 0),
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
	sp.Me = false
	sp.DeckWeight = 0
	deckhid := Deck{}
	for _, _ = range sp.Deck {
		deckhid.Add(&Card{})
	}
	sp.Deck = deckhid
}

func NewState(g *Game, p *Player) *State {
	stateplayers := make([]StatePlayer, 0)
	var sp StatePlayer
	for _, player := range g.Players {
		playing := g.GetCurrentPlayer().Id == player.Id
		sp = NewStatePlayer(player, playing)
		if sp.Id != p.Id {
			sp.HideInfos()
		}
		stateplayers = append(stateplayers, sp)
	}
	return &State{
		PlayDeck: *g.PlayDeck,
		LastLog:  *g.LastLog,
		Round:    g.Round,
		Players:  stateplayers,
		Started:  g.Started,
	}
}

func (s State) Json() (string, []byte) {
	data, _ := json.Marshal(s)
	return string(data), data
}
