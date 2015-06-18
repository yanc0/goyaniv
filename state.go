package main

import (
  "encoding/json"
)

type State struct {
  PlayDeck Deck `json:"playdeck"`
  PlayersScore map[string]int `json:"playersscore"`
  PlayersNumCard map[string]int `json:"playersnumcard"`
  PlayerTurn string `json:"playerturn"`
  LogLine string `json:"logline"`
  Round int `json:"round"`
  Launched bool `json:"launched"`
}

func (s State) Json() (string, []byte) {
  data, _ := json.Marshal(s)
  return string(data), data
}
