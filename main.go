package main

import (
  "fmt"
  "github.com/olahol/melody"
  "github.com/gin-gonic/gin"
  "net/http"
)

type Goyaniv struct {
  Ws *melody.Melody
  Games []*Game
}

func (goya *Goyaniv) AddGame (game *Game) {
  goya.Games = append(goya.Games, game)
}

func (goya *Goyaniv) FindGame(s *melody.Session) (*Game){
  for _, game := range goya.Games {
    for _, player := range game.Players {
      if s == player.Session {
        return game
      }
    }
  }
  return nil
}

func main() {
  d := Deck(make([]Card, 0))
  d.Init()
  d.Shuffle()

  r := gin.Default()
  m := melody.New()

  game := Game{"Partie Yaniv", make([]*Player, 0), &d, 0}
  goya := Goyaniv{m, make([]*Game, 0)}
  goya.AddGame(&game)

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
    fmt.Println("Handle Connect, new player")
    deck := Deck{}
    player := Player{"Guest", s, &deck, "playing"}
    game.Players = append(game.Players, &player)
  })
  r.Run(":5000") 
}
