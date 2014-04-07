package main

import "math/rand"
import "time"
import "fmt"

func main() {
	rand.Seed(time.Now().Unix())

	deck := new(Deck)
	deck.Init()

	for !deck.IsEmpty() {
		fmt.Println(deck.Draw())
	}

	fmt.Println("")

	deck.Init()

	for !deck.IsEmpty() {
		fmt.Println(deck.Draw())
	}
}