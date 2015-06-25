package goyaniv

import (
	"github.com/olahol/melody"
)

type Player struct {
	Name      string          `json:"name"`
	Session   *melody.Session `json:"session"`
	Deck      *Deck           `json:"session"`
	State     string          `json:"state"`
	Score     int             `json:"score"`
	Connected bool            `json:"connected"`
	Id        string          `json:"id"`
	Key       string          `json:"key"`
}
