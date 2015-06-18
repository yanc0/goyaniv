package main

type Game struct {
  Name string `json:"name"`
  Players []*Player `json:"players"`
  Playing *Player `json:"playing"`
  Fastplaying *Player `json:"playing"`
  MiddleDeck *Deck `json:"middledeck"`
  PlayDeck *Deck `json:"playdeck"`
  Round int `json:"round"`
  Url string `json:"url"`
  Launched bool `json:"launched"`
}

func (g *Game) GetState() (State){
  s := State{PlayDeck: make([]Card, 0),
             PlayersScore: make(map[string]int,0),
             PlayersNumCard: make(map[string]int,0)}
  if g.PlayDeck != nil {
    for _, card := range *g.PlayDeck {
      s.PlayDeck.Add(card)
    }
  }
  for _, player := range g.Players {
    s.PlayersScore[player.Name] = player.Score
    s.PlayersNumCard[player.Name] = len(*player.Deck)
    if g.Playing == player {
      s.PlayerTurn = player.Name
    }
  }
  return s
}
