package main

import "math/rand"
import "time"
import "os"

var players []*Player

func main() {
	rand.Seed(time.Now().Unix())

	players = make([]*Player, len(os.Args)-1)
	var startPlayer = 0
	

	for i, _ := range players {
		players[i] = new(Player)
		players[i].Init(os.Args[i+1], i, len(players), startPlayer)
	}

	deck := new(Deck)
	deck.Init()

	for i, _ := range players {
		players[i].Send("draw", deck.Draw())
		players[i].Forfeit()
	}
}