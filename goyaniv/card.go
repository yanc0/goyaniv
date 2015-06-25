package goyaniv

import (
	"encoding/json"
	"log"
)

type Card struct {
	Id     int    `json:"id"`
	Value  int    `json:"value"`  // 0=JKR, 1=AS, 11=J, 12=Q, 13=K
	Symbol string `json:"symbol"` // spade, heart, diamond, club
}

func (c *Card) Color() string {
	if c.Symbol == "spade" || c.Symbol == "club" {
		return "black"
	}
	return "red"
}

func (c *Card) Weight() int {
	if c.Value >= 10 {
		return 10
	}
	return c.Value
}

func (c *Card) Json() (string, []byte) {
	data, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	return string(data), data
}
