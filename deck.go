package main

import "math/rand"

type Deck struct{
	cards []Card
}

func (d *Deck) Init() {
	//load the cards into the deck
	d.cards = []Card{
		Soldier, Soldier, Soldier, Soldier, Soldier,
		Clown, Clown,
		Knight, Knight,
		Priestess, Priestess,
		Wizard, Wizard,
		General,
		Minister,
		Princess,
	}
	
	//shuffle the cards
	perm := rand.Perm(len(d.cards))
	var newcards = make([]Card, 0, len(d.cards))
	for _, i := range perm {
		newcards = append(newcards, d.cards[i])
	}
	d.cards = newcards
}

func (d *Deck) Draw() Card {
	n := len(d.cards) - 1
	drawnCard := d.cards[n]
	
	d.cards = d.cards[0:n]
	return drawnCard
}

func (d *Deck) IsEmpty() bool {
	return (len(d.cards) == 0)
}