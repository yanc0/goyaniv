package main

import (
  "encoding/json"
  "github.com/olahol/melody"
  "fmt"
)

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

func ActionPut(goya *Goyaniv, s *melody.Session, action *Action) {
  game := goya.FindGameWithSession(s)
  for _, player := range game.Players {
    if player.Session == s {
      decktmp := Deck{}
      for _, id := range action.PutCards {
        decktmp.Add(player.Deck.TakeCardID(id))
      }
      for _, card := range decktmp {
        game.PlayDeck.Add(card)
      }
      cardtaken := game.PlayDeck.TakeCardID(action.TakeCard)
      if cardtaken == nil {
        fmt.Println("Card does not exist in deck")
      }
      player.Deck.Add(cardtaken)
    }
  }
}

type Action struct {
  Name string `json:"name"`
  PutCards []int `json:"putcards"`
  TakeCard int `json:"takecard"`
  Option string `json:"option"`
}

func FireAction (goya *Goyaniv, s *melody.Session, jsonrcv []byte) {
  action := &Action{}
  json.Unmarshal(jsonrcv, &action)
  fmt.Println(*action)
  switch action.Name {
    case "draws":
      ActionDraws(goya, s)
      BroadcastState(goya, s)
    case "yaniv":
      ActionYaniv(goya, s)
      BroadcastState(goya, s)
    case "score":
      ActionGetScore(goya, s)
    case "mydeck":
      ActionMyDeck(goya, s)
    case "set":
      ActionSetName(goya, s, action.Option)
      BroadcastState(goya, s)
    case "put":
      ActionPut(goya, s, action)
      BroadcastState(goya, s)
  }
}
