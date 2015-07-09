package goyaniv

import (
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"net/http"
)

type Server struct {
	Ws     *melody.Melody
	Routes *gin.Engine
	Games  []*Game
}

func New() *Server {
	ws := melody.New()
	routes := gin.Default()
	games := make([]*Game, 0)

	return &Server{
		Ws:     ws,
		Routes: routes,
		Games:  games,
	}
}

func (s *Server) AddGame(game *Game) {
	s.Games = append(s.Games, game)
}

func (s *Server) FindGameWithPlayerId(id string) *Game {
	for _, game := range s.Games {
		for _, player := range game.Players {
			if id == player.Id {
				return game
			}
		}
	}
	return nil
}

func (s *Server) GetGameWithURL(url string) *Game {
	for _, game := range s.Games {
		if game.Url == url {
			return game
		}
	}
	return nil
}

func (s *Server) FindPlayerWithId(id string) *Player {
	for _, game := range s.Games {
		for _, player := range game.Players {
			if player.Id == id {
				return player
			}
		}
	}
	return nil
}

func (srv *Server) RoutesInit() {
	srv.Routes.GET("/jq.js", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "html/jq.js")
	})

	srv.Routes.GET("/goyaniv.js", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "html/goyaniv.js")
	})
	srv.Routes.GET("/game/:name", func(c *gin.Context) {
		cookiekey, _ := c.Request.Cookie("goyanivkey")
		cookieid, _ := c.Request.Cookie("goyanivid")
		if cookieid == nil || cookiekey == nil {
			cookieid := CreateCookie("goyanivkey", GenerateUnique())
			cookiekey := CreateCookie("goyanivid", GenerateUnique())
			http.SetCookie(c.Writer, cookieid)
			http.SetCookie(c.Writer, cookiekey)
		}
		http.ServeFile(c.Writer, c.Request, "html/game2.html")
	})

	srv.Routes.GET("/game/:name/ws", func(c *gin.Context) {
		srv.Ws.HandleRequest(c.Writer, c.Request)
	})

	srv.Ws.HandleMessage(func(s *melody.Session, msg []byte) {
		FireMessage(srv, s, msg)
	})

	srv.Ws.HandleDisconnect(func(s *melody.Session) {
		FireDisconnect(srv, s)
	})

	srv.Ws.HandleConnect(func(s *melody.Session) {
		FireConnect(srv, s)
	})
}

func (s *Server) Run() {
	s.RoutesInit()
	s.Routes.Run(":5000")
}
