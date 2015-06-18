package main

import (
  "fmt"
  "github.com/olahol/melody"
  "github.com/gin-gonic/gin"
  "net/http"
  "github.com/satori/go.uuid"
)

type Goyaniv struct {
  Ws *melody.Melody
  Games []*Game
}

func (goya *Goyaniv) AddGame (game *Game) {
  goya.Games = append(goya.Games, game)
}

func (goya *Goyaniv) FindGameWithSession(s *melody.Session) (*Game){
  for _, game := range goya.Games {
    for _, player := range game.Players {
      if s == player.Session {
        return game
      }
    }
  }
  return nil
}

func (goya *Goyaniv) FindGameWithUrl(url string) (*Game){
  for _, game := range goya.Games {
    if game.Url == url {
      return game
    }
  }
  return nil
}

func main() {

  r := gin.Default()
  m := melody.New()

  goya := Goyaniv{m, make([]*Game, 0)}

  r.GET("/game/:name", func(c *gin.Context) {
      http.ServeFile(c.Writer, c.Request, "game.html")
  })

  r.GET("/game/:name/ws", func(c *gin.Context) {
      m.HandleRequest(c.Writer, c.Request)
  })

  m.HandleMessage(func(s *melody.Session, msg []byte) {
    FireAction(&goya, s, msg)
    fmt.Println(string(msg))
  })

  m.HandleConnect(func(s *melody.Session) {
    fmt.Print("Handle Connect, new player: ")
    game := goya.FindGameWithUrl(s.Request.URL.Path)
    if (game == nil) {
      middledeck := Deck{}
      middledeck.Init()
      middledeck.Shuffle()
      playdeck := Deck{}
      playdeck.Add(middledeck.TakeCard())
      game = &Game{Name: "Partie Yaniv",
                  Players: make([]*Player, 0),
                  MiddleDeck: &middledeck,
                  PlayDeck: &playdeck,
                  Round: 0,
                  Url: s.Request.URL.Path}
      goya.AddGame(game)
    }
    deck := Deck{}
    for i := 0; i < 5; i++ {
      deck.Add(game.MiddleDeck.TakeCard())
    }
    uniquename := uuid.NewV4().String()
    fmt.Println(uniquename)
    player := Player{Name: uniquename, Session: s, Deck: &deck, State:"playing"}
    game.Players = append(game.Players, &player)
  })
  r.Run(":5000") 
}
