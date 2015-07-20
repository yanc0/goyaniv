package goyaniv

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

type Deck []*Card

func (d Deck) Score() int {
	var i int
	for _, card := range d {
		i += card.Weight()
	}
	return i
}

func NewCompleteDeck() *Deck {
	deck := &Deck{}
	symbols := make([]string, 4)
	symbols[0] = "spade"
	symbols[1] = "heart"
	symbols[2] = "diam"
	symbols[3] = "club"
	var j int = 1
	for _, symbol := range symbols {
		for i := 1; i < 14; i++ {
			card := &Card{j, i, symbol}
			deck.Add(card)
			j++
		}
	}
	return deck
}
func (d *Deck) AddDeck(deckadd *Deck) {
	for _, _ = range *deckadd {
		d.Add(deckadd.TakeCard())
	}
}

func (d *Deck) Weight() int {
	var weight int
	for _, card := range *d {
		weight += card.Weight()
	}
	return weight
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range *d {
		j := rand.Intn(i + 1)
		(*d)[i], (*d)[j] = (*d)[j], (*d)[i]
	}
}

func (d *Deck) Add(c *Card) {
	*d = append(*d, c)
}

func (d *Deck) TakeCard() *Card {
	//shift first card
	c := (*d)[0]
	*d = (*d)[1:]
	return c
}

func (d Deck) Len() int {
	return len(d)
}

func (d *Deck) IsValid() bool {
	return d.IsSequence() || d.IsMultiple() || d.Len() == 1
}

func (d Deck) IsMultiple() bool {
	if d.Len() < 2 {
		return false
	}
	value := d[0].Value
	for _, card := range d {
		if card.Value != value {
			return false
		}
	}
	return true
}

func (d Deck) IsSequence() bool {
	dtmp := d
	var min, max int = 14, 0
	k := 0
	// If less than 3 cards, not a seq.
	if len(dtmp) < 3 {
		return false
	}
	// check if all cards have same Symbol
	Symbol := ""
	for _, card := range dtmp {
		if Symbol == "" && card.Value > 0 {
			Symbol = card.Symbol
		}
		if card.Value > 0 && card.Symbol != Symbol {
			return false
		}
	}
	// if there is more than double card, not a seq
	for i := range dtmp {
		for j := range d {
			if dtmp[i] == d[j] {
				k++
			}
		}
		if k > 1 {
			return false
		} else {
			k = 0
		}
	}
	// Get max Value and min Value
	for _, card := range dtmp {
		if card.Value < min && card.Value != 0 {
			min = card.Value
		}
		if card.Value > max {
			max = card.Value
		}
	}
	// final check
	if max-min <= len(dtmp)-1 {
		return true
	}
	return false
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Deck) Less(i, j int) bool {
	return d[i].Value < d[j].Value
}

func (d Deck) Json() string {
	data, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (d *Deck) TakeCardID(id int) *Card {
	if d == nil {
		return nil
	}
	for i, card := range *d {
		if card.Id == id {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return card
		}
	}
	return nil
}
