package main

import "fmt"
//import "sort"
import "time"
import "math/rand"

type Card struct {
  value int // 0=JKR, 1=AS, 11=J, 12=Q, 13=K
  symbol string // spade, heart, diamond, club
}

type Deck []Card

func (c *Card) Color() string {
  if c.symbol == "spade" || c.symbol == "club" {
    return "black"
  }
  return "red"
}

func (c *Card) Weight() int {
  if c.value >= 10 {
    return 10
  }
  return c.value
}


func (d *Deck) Init(){
  symbols := make([]string, 4)
  symbols[0] = "spade"
  symbols[1] = "diamond"
  symbols[2] = "heart"
  symbols[3] = "club"
  for i := 0; i < 14; i++ {
    for _, symbol := range symbols {
      d.Add(Card{i, symbol})
    }
  }
}

func (d Deck) Shuffle() {
  rand.Seed(time.Now().UTC().UnixNano())
  for i := range d {
    j := rand.Intn(i + 1)
    d[i], d[j] = d[j], d[i]
  }
}

func (d *Deck) Add (c Card) {
  *d = append(*d, c)
}

func (d Deck) TakeCard() Card {
  //shift first card
  c := Card{}
  c, d = d[0], d[1:]
  return c
}

func (d Deck) Len() int {
  return len(d)
}

func (d Deck) IsSequence() bool {
  dtmp := d
  var min, max int = 14, 0
  k := 0
  // If less than 3 cards, not a seq.
  if len(dtmp) < 3 {
    return false
  }
  // check if all cards have same symbol
  symbol := ""
  for _, card := range dtmp {
    if symbol == "" && card.value > 0 {
      symbol = card.symbol
    }
    if card.value > 0 && card.symbol != symbol {
      return false
    }
  }
  // if there is more than double card, not a seq
  for i := range dtmp {
    for j := range d {
      if dtmp[i] == d[j]{
        k++
      }
    }
    if k > 1 {
      return false
    } else {
      k = 0
    }
  }
  // Get max value and min value
  for _, card := range dtmp {
    if card.value < min && card.value != 0 {
      min = card.value
    }
    if card.value > max {
      max = card.value
    }
  }
  // final check
  if max - min <= len(dtmp) - 1 {
    return true
  }
  return false
}

func (d Deck) Swap(i, j int) {
  d[i], d[j] = d[j], d[i]
}

func (d Deck) Less(i, j int) bool {
  return d[i].value < d[j].value
}

func main() {
  d := Deck(make([]Card, 0))
  d.Add(Card{4, "spade"})
  d.Add(Card{5, "heart"})
  d.Add(Card{6, "spade"})
  fmt.Println(d.IsSequence())
}
