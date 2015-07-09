package goyaniv

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
)

func UnicastState(game *Game, player *Player, error string) {
	state := NewState(game, player)
	state.Error = error
	if player.Session != nil {
		_, data := state.Json()
		fmt.Println("Sending state to", player.Id)
		player.Session.Write(data)
	}
}

func BroadcastState(game *Game) {
	for _, player := range game.Players {
		UnicastState(game, player, "")
	}
}

func ActionSetName(p *Player, name string) {
	p.Name = name
}

func ActionYaniv(game *Game, p *Player) string {
	if p.Id == game.GetCurrentPlayer().Id {
		if p.Deck.Weight() <= game.YanivAt {
			for _, player := range game.Players {
				if p.Deck.Weight() > game.YanivAt {
					player.WantsAsaf = "no"
				} else {
					player.WantsAsaf = "noanswer"
				}
			}
			p.Yaniv = true
			p.WantsAsaf = "yes"
		} else {
			error := "You cant' yaniv yet"
			UnicastState(game, p, error)
			return error
		}
	} else {
		error := "It is not your turn"
		UnicastState(game, p, error)
		return error
	}
	return "noerror"
}

func ActionAsaf(game *Game, player *Player, answer string) string {
	if answer != "yes" || answer != "no" {
		error := "Bad answer"
		UnicastState(game, player, error)
		return error
	}
	if len(game.PlayersWantsAsaf()) == 0 {
		error := "Nobody yaniv yet, you can't asaf"
		UnicastState(game, player, error)
		return error
	} else {
		if player.WantsAsaf == "noanswer" {
			player.WantsAsaf = answer
			if answer == "yes" {
				player.Asaf = game.GetAsafRank()
			}
		}
	}
	return "noerror"
}

func ActionPut(game *Game, p *Player, action *Action) (err string) {
	if len(game.PlayersWantsAsaf()) == 0 {
		error := "Game stopped, you can asaf only"
		UnicastState(game, p, error)
		return error
	}
	if p.Id == game.GetCurrentPlayer().Id {
		decktmp := Deck{}
		for _, id := range action.PutCards {
			c := p.Deck.TakeCardID(id)
			if c == nil {
				for _, card := range decktmp {
					p.Deck.Add(card)
				}
				error := "Put cards does not exists in player deck"
				UnicastState(game, p, error)
				return error
			}
			decktmp.Add(c)
		}
		if !decktmp.IsValid() {
			for _, card := range decktmp {
				p.Deck.Add(card)
			}
			error := "invalid combination"
			UnicastState(game, p, error)
			return error
		}
		if action.TakeCard == 0 {
			cardtaken := game.MiddleDeck.TakeCard()
			if cardtaken == nil {
				error := "Card does not exist in deck"
				UnicastState(game, p, error)
				return error
			}
			p.Deck.Add(cardtaken)
		} else {
			cardtaken := game.PlayDeck.TakeCardID(action.TakeCard)
			if cardtaken == nil {
				error := "Card does not exist in deck"
				UnicastState(game, p, error)
				return error
			}
			p.Deck.Add(cardtaken)
		}
		game.TrashDeck.AddDeck(game.PlayDeck)
		for _, card := range decktmp {
			game.PlayDeck.Add(card)
		}
		return "noerror"
	} else if p == game.GetFastPlayer() {
		return "fastplayer"
	}
	error := "It is not your turn"
	UnicastState(game, p, error)
	return error
}

type Action struct {
	Name     string `json:"name"`
	PutCards []int  `json:"putcards"`
	TakeCard int    `json:"takecard"`
	Option   string `json:"option"`
}

func (a *Action) isValid() bool {
	if a == nil {
		return false
	}
	if a.Name == "put" && len(a.PutCards) == 0 {
		return false
	}
	return true
}

func JSONToAction(jsn []byte) *Action {
	action := &Action{}
	json.Unmarshal(jsn, &action)
	return action
}

func GetGameWithSession(s *melody.Session) *Game {
	game := &Game{}
	return game
}

func GetPlayerWithSession(s *melody.Session) *Player {
	player := &Player{}
	return player
}

func FireConnect(srv *Server, s *melody.Session) {
	gameUrl := s.Request.URL.Path
	cookiekey, _ := s.Request.Cookie("goyanivkey")
	playerkey := cookiekey.Value
	cookieid, _ := s.Request.Cookie("goyanivid")
	playerid := cookieid.Value
	game := srv.GetGameWithURL(gameUrl)
	// create game if it does not exists
	if game == nil {
		game = NewGame(gameUrl)
		srv.AddGame(game)
	}
	player := game.GetPlayer(playerid, playerkey)
	if player == nil {
		playerdeck := &Deck{}
		for i := 0; i < 5; i++ {
			playerdeck.Add(game.MiddleDeck.TakeCard())
		}
		player = &Player{
			Name:    GenerateUnique(),
			Session: s,
			Id:      playerid,
			Deck:    playerdeck,
			Key:     playerkey,
		}
		game.AddPlayer(player)
		fmt.Println("Player ID", player.Id, "connected")
	} else {
		player.Session = s
		player.Connected = true
		fmt.Println("Player ID", player.Id, "reconnected")
	}
	player.Connected = true
	BroadcastState(game)
}

func FireDisconnect(srv *Server, s *melody.Session) {
	game := srv.GetGameWithURL(s.Request.URL.Path)
	cookieid, _ := s.Request.Cookie("goyanivid")
	cookiekey, _ := s.Request.Cookie("goyanivkey")
	player := game.GetPlayer(cookieid.Value, cookiekey.Value)
	player.Session = nil
	player.Connected = false
	fmt.Println(player.Id, "Disconnected")
	BroadcastState(game)
}

func FireMessage(srv *Server, s *melody.Session, jsn []byte) bool {
	action := JSONToAction(jsn)
	game := srv.GetGameWithURL(s.Request.URL.Path)
	cookieid, _ := s.Request.Cookie("goyanivid")
	cookiekey, _ := s.Request.Cookie("goyanivkey")
	player := game.GetPlayer(cookieid.Value, cookiekey.Value)
	if player.Key != cookiekey.Value && player.Id != cookieid.Value {
		return false
	}
	fmt.Println(action, "player", player.Id)

	if !action.isValid() {
		UnicastState(game, player, "Action not valid")
		return false
	}

	if game == nil || player == nil {
		return false
	}

	switch action.Name {
	case "name":
		ActionSetName(player, action.Option)
		BroadcastState(game)
	case "put":
		err := ActionPut(game, player, action)
		if err == "noerror" {
			game.NextPlayer()
			BroadcastState(game)
		} else {
			fmt.Println(err)
			return false
		}
	case "yaniv":
		err := ActionYaniv(game, player)
		if err == "noerror" {
			fmt.Println(player.Name, "just yanived !")
			game.UpdateScores()
			BroadcastState(game)
		} else {
			fmt.Println(err)
		}
	case "asaf":
		err := ActionAsaf(game, player, action.Option)
		if err == "noerror" {
			fmt.Println(player.Name, "Asafed")
			game.UpdateScores()
			BroadcastState(game)
		} else {
			fmt.Println(err)
		}
	}
	return true
}
