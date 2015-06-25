package goyaniv

import (
	"encoding/json"
)

type State struct {
	PlayDeck       Deck              `json:"playdeck"`
	PlayersScore   map[string]int    `json:"playersscore"`
	PlayersNumCard map[string]int    `json:"playersnumcard"`
	PlayersName    map[string]string `json:"playersname"`
	PlayerTurn     string            `json:"playerturn"`
	LogLine        string            `json:"logline"`
	Round          int               `json:"round"`
	Launched       bool              `json:"launched"`
	Myself         Player            `json:"myself"`
	MyDeck         Deck              `json:"mydeck"`
}

func (s State) Json() (string, []byte) {
	data, _ := json.Marshal(s)
	return string(data), data
}
