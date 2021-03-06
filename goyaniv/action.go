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

func ActionReady(game *Game, p *Player, option string) string {
	if option == "no" {
		if game.Started {
			return "Game already started, you have to stay ready"
		}
		p.Ready = false
	} else if option == "yes" {
		p.Ready = true
	} else {
		UnicastState(game, p, "Invalid Ready option")
	}

	for _, player := range game.Players {
		if player.Ready == false {
			return "wait"
		}
	}
	return "start"
}

func ActionYaniv(game *Game, p *Player) string {
	if p.Id == game.GetCurrentPlayer().Id {
		if p.Deck.Weight() <= game.YanivAt {
			for _, player := range game.Players {
				if player.Deck.Weight() > game.YanivAt || player.IsSpectator() {
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
	if !game.Started {
		return "Game have not started yet"
	}
	for _, id := range action.PutCards {
		putcard := p.Deck.TakeCardID(id)
		if putcard == nil {
			// player want to send card he does not have
			// put back taken card in is deck
			p.Deck.AddDeck(decktmp)
			return "Put cards does not exists in player deck"
		}
		decktmp.Add(putcard)
	}
	if decktmp.Len() == 1 || decktmp.IsMultiple() || decktmp.IsSequence() {
		// Classic case
		if p == game.GetCurrentPlayer() {
			if action.TakeCard == 0 {
				game.TrashDeck.AddDeck(game.PlayDeck)
				game.PlayDeck = decktmp
				p.Deck.Add(game.MiddleDeck.TakeCard())
				return "noerror"
			} else {
				takecard := game.PlayDeck.TakeCardID(action.TakeCard)
				if takecard == nil {
					p.Deck.AddDeck(decktmp)
					return "Card does not exists in playdeck"
				} else {
					playdeckcopy := *game.PlayDeck
					decktmpcopy := *decktmp
					deckverif := &Deck{}
					deckverif.AddDeck(&playdeckcopy)
					deckverif.AddDeck(&decktmpcopy)
					deckverif.Add(takecard)
					if deckverif.IsMultiple() && deckverif.Len() == 4 {
						game.PlayDeck = deckverif
						return "fastplay"
					} else {
						game.TrashDeck.AddDeck(game.PlayDeck)
						game.PlayDeck = decktmp
						p.Deck.Add(takecard)
						return "noerror"
					}
				}
			}
		} else {
			playdeckcopy := *game.PlayDeck
			decktmpcopy := *decktmp
			decktmp.AddDeck(&playdeckcopy)
			// if it's not your turn you can only put multiple
			// and with two condition:
			// - you can fast play
			// - you complete the four multiple
			if decktmp.IsMultiple() && game.LastLog.Action == "put" {
				if (game.GetFastPlayer() == p && game.LastLog.TakeCard.Id == 0) || decktmp.Len() == 4 {
					game.PlayDeck = decktmp
					return "fastplay"
				} else {
					p.Deck.AddDeck(&decktmpcopy)
				}
			} else {
				// else it's not your turn
				p.Deck.AddDeck(&decktmpcopy)
			}
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
		state := "playing"
		ready := false

		if !game.Started || len(game.Players) >= 5 {
			for i := 0; i < 5; i++ {
				playerdeck.Add(game.MiddleDeck.TakeCard())
			}
		} else {
			state = "spectator"
			ready = true
		}
		player = &Player{
			Name:    GenerateUnique(),
			Session: s,
			Id:      playerid,
			Deck:    playerdeck,
			Key:     playerkey,
			State:   state,
			Ready:   ready,
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
	if !game.Started {
		game.DeletePlayer(player.Id)
		game.debug_is_game_consistent()
	}
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
	case "ready":
		ret := ActionReady(game, player, action.Option)
		fmt.Println("game ", ret)
		if ret == "start" {
			game.Launch()
		}
		BroadcastState(game)
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
