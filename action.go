package main

import (
  "encoding/json"
  "github.com/olahol/melody"
)

func ActionTakeCard(goya *Goyaniv, s *melody.Session) {
  game := *goya.FindGame(s)
  c := game.Middledeck.TakeCard()
  _, data := c.Json()
  s.Write(data)
}

func ActionGetScore(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGame(s)
  for _, player := range (*game).Players{
    if player.Session == s {
      data, _ := json.Marshal(player.Deck.Score())
      s.Write(data)
    }
  }
}

func ActionYaniv(goya *Goyaniv, s *melody.Session) {
  game := goya.FindGame(s)
  for _, player := range game.Players {
    player.Session.Write([]byte("yaniv"))
  }
}

func FireAction (goya *Goyaniv, s *melody.Session, action []byte) {
  a := string(action)
  switch a {
    case "takecard":
      ActionTakeCard(goya, s)
    case "yaniv":
      ActionYaniv(goya, s)
    case "score":
      ActionGetScore(goya, s)
  }
}
