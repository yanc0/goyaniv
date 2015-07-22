package goyaniv

import (
	"encoding/json"
	"fmt"
	"github.com/olahol/melody"
)

type Action struct {
	Name       string `json:"name"`
	PutCards   []int  `json:"putcards"`
	TakeCard   int    `json:"takecard"`
	Option     string `json:"option"`
	PlayerName string `json:"playername"`
}

func UnicastState(game *Game, player *Player, error string) {
	state := NewState(game, player)
	state.Error = error
	if player.Session != nil {
		_, data := state.Json()
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
				if player.Deck.Weight() > game.YanivAt {
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
	if answer == "yes" || answer == "no" {
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
	} else {
		error := "Bad answer"
		UnicastState(game, player, error)
		return error
	}
	return "noerror"
}

func ActionPut(game *Game, p *Player, action *Action) (err string) {
	decktmp := &Deck{}
	if len(game.PlayersWantsAsaf()) != 0 {
		return "Game stopped, you can asaf only"
	}
	for _, id := range action.PutCards {
		putcard := p.Deck.TakeCardID(id)
		if putcard == nil {
			// player want to send card he does not have
			// put back taken card in is deck
			p.Deck.AddDeck(decktmp)
			return "Put cards does not exists in player deck"
		}
		fmt.Println(putcard)
		decktmp.Add(putcard)
	}
	if decktmp.Len() == 1 || decktmp.IsMultiple() || decktmp.IsSequence() {
		// Classic case
		if p == game.GetCurrentPlayer() {
			if action.TakeCard == 0 {
				game.PlayDeck = decktmp
				p.Deck.Add(game.MiddleDeck.TakeCard())
				return "noerror"
			} else {
				takecard := game.PlayDeck.TakeCardID(action.TakeCard)
				if takecard == nil {
					p.Deck.AddDeck(decktmp)
					return "Card does not exists in playdeck"
				} else {
					game.TrashDeck.AddDeck(game.PlayDeck)
					game.PlayDeck = decktmp
					p.Deck.Add(takecard)
					return "noerror"
				}
			}
		} else {
			playdeckcopy := *game.PlayDeck
			decktmp.AddDeck(&playdeckcopy)
			// if it's not your turn you can only put multiple
			// and with two condition:
			// - you can fast play
			// - you complete the four multiple
			if decktmp.IsMultiple() {
				if game.GetFastPlayer() == p || decktmp.Len() == 4 {
					game.TrashDeck.AddDeck(game.PlayDeck)
					game.PlayDeck = decktmp
					return "fastplay"
				}
			}
			// else it's not your turn
			p.Deck.AddDeck(decktmp)
			return "It is not your turn"
		}
	}
	p.Deck.AddDeck(decktmp)
	return "You put invalid card combination"
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
			State:   "playing",
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

	if !action.isValid() {
		UnicastState(game, player, "Action not valid")
		return false
	}
	action.PlayerName = player.Name
	fmt.Println(action)
	switch action.Name {
	case "name":
		ActionSetName(player, action.Option)
		action.ToLog(game)
		BroadcastState(game)
	case "put":
		err := ActionPut(game, player, action)
		if err == "noerror" {
			game.NextPlayer()
			action.ToLog(game)
			BroadcastState(game)
		} else if err == "fastplay" {
			action.ToLog(game)
			BroadcastState(game)
		} else {
			UnicastState(game, player, err)
			return false
		}
	case "yaniv":
		err := ActionYaniv(game, player)
		if err == "noerror" {
			action.ToLog(game)
			game.UpdateScores()
			BroadcastState(game)
		} else {
			fmt.Println(err)
		}
	case "asaf":
		err := ActionAsaf(game, player, action.Option)
		if err == "noerror" {
			fmt.Println(player.Name, "Asafed")
			action.ToLog(game)
			game.UpdateScores()
			BroadcastState(game)
		} else {
			fmt.Println(err)
		}
	}
	return true
}
