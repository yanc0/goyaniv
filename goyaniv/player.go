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
	Ready     bool            `json:"ready"`
	Id        string          `json:"id"`
	Key       string          `json:"key"`
	WantsAsaf string          `json:"wantsasaf"`
	Asaf      int             `json:"asaf"`
	Yaniv     bool            `json:"yaniv"`
}

type ListPlayer []*Player

func (lp ListPlayer) Len() int {
	return len(lp)
}

func (lp ListPlayer) Swap(i, j int) {
	lp[i], lp[j] = lp[j], lp[i]

}
func (lp ListPlayer) Less(i, j int) bool {
	//return len(lp[i]) < len(lp[j])
	if lp[i].Deck.Weight() < lp[j].Deck.Weight() {
		return true
	}
	if lp[i].Deck.Weight() == lp[j].Deck.Weight() {
		if lp[j].Yaniv {
			return true
		}
		return lp[i].Asaf < lp[j].Asaf
	}
	return false
}
