package main

import "math/rand"
import "time"
import "os"

var players [4]*Player

func main() {
	rand.Seed(time.Now().Unix())

	n := len(os.Args)
	if n != 5 {
		panic("Not enough arguments")
	}
	
	

	for i, _ := range players {
		players[i] = new(Player)
		players[i].Init(os.Args[i+1], i)
	}

	deck := new(Deck)
	deck.Init()

	for i, _ := range players {
		players[i].Send("draw", deck.Draw())
		players[i].Forfeit()
	}
}