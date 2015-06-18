package main

type Game struct {
  Name string `json:"name"`
  Players []*Player `json:"players"`
  Middledeck *Deck `json:"middledeck"`
  Round int `json:"round"`
}
