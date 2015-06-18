package main

import (
  "encoding/json"
  "strings"
  "github.com/olahol/melody"
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

func ActionDraws(goya *Goyaniv, s *melody.Session) {
  game := *goya.FindGameWithSession(s)
  c := game.MiddleDeck.TakeCard()
  _, data := c.Json()
  s.Write(data)
}

func ActionGetScore(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGameWithSession(s)
  for _, player := range (*game).Players{
    if player.Session == s {
      data, _ := json.Marshal(player.Deck.Score())
      s.Write(data)
    }
  }
}

func ActionMyDeck(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGameWithSession(s)
  for _, player := range (*game).Players{
    if player.Session == s {
      data, _ := json.Marshal(player.Deck)
      s.Write(data)
    }
  }
}

func ActionYaniv(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGameWithSession(s)
  for _, player := range game.Players {
    player.Session.Write([]byte("yaniv"))
  }
}

func BroadcastState(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGameWithSession(s)
  _, data := game.GetState().Json()
  for _, player := range game.Players {
    player.Session.Write(data)
  }
}

func ActionSetName(goya *Goyaniv, s *melody.Session, name string) {
  game := goya.FindGameWithSession(s)
  for _, player := range game.Players {
    if player.Session == s {
      player.Name = name
    }
  }
}

func FireAction (goya *Goyaniv, s *melody.Session, action []byte) {
  a := string(action)
  switch strings.Split(a, " ")[0] {
    case "draws":
      ActionDraws(goya, s)
    case "yaniv":
      ActionYaniv(goya, s)
    case "score":
      ActionGetScore(goya, s)
    case "mydeck":
      ActionMyDeck(goya, s)
    case "set":
      ActionSetName(goya, s, strings.Split(a, " ")[1])
  }
  BroadcastState(goya, s)
}
