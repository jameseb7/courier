package main

import "math/rand"
import "time"
import "os"
import "fmt"
import "strings"

var players []*Player

func main() {
	rand.Seed(time.Now().Unix())

	players = make([]*Player, len(os.Args)-1)
	var startPlayer = 0
	
	var deck = new(Deck)
	
	for i, _ := range players {
			players[i] = new(Player)
	}

	for won, _ := PlayerWon(); !won; {
		deck.Init()
		for i, _ := range players {
			players[i].Init(os.Args[i+1], i, len(players), startPlayer)
			players[i].Send("draw", deck.Draw())
		}

		_ = deck.Draw() //throw away a card to prevent perfect information

		var currentPlayer = startPlayer
		for !deck.IsEmpty() && (SurvivingPlayers() > 1) {
			var fullMove, action, cardStr, targetCardStr string
			var target, n  int
			var card, cardDrawn Card

			if players[currentPlayer].lost {
				goto endTurn
			}

			SendAll("player", currentPlayer)
			
			cardDrawn = deck.Draw()
			players[currentPlayer].Send("draw", cardDrawn)
			players[currentPlayer].AddToHand(cardDrawn)
			

			fullMove = players[currentPlayer].Recieve()
			
			
			n, _ = fmt.Sscan(fullMove, &action, &cardStr, &target, &targetCardStr)
			
			if n < 2 {
				//malformed move
				fmt.Println("ERROR: Player", currentPlayer, "made an illegal play")
				players[currentPlayer].Forfeit()
				goto endTurn
			}

			if !strings.EqualFold(action, "play") {
				//must be "forfeit" or some invalid action
				fmt.Println("ERROR: Player", currentPlayer, "made an illegal play")
				players[currentPlayer].Forfeit()
				goto endTurn
			}

			card = ParseCard(cardStr)
			if players[currentPlayer].RemoveFromHand(card) {
				//trying to play a card you don't have is illegal
				fmt.Println("ERROR: Player", currentPlayer, "played a noexistent card")
				players[currentPlayer].Forfeit()
				goto endTurn
			}

			switch card {
			case Princess:
				SendAll("played", currentPlayer, card)
				players[currentPlayer].Forfeit()
				SendAll("out", currentPlayer, players[currentPlayer].HandString())
			}

		endTurn:
			currentPlayer++
			if currentPlayer >= len(players) {
				currentPlayer = 0
			}
		}
		startPlayer++
		if startPlayer >= len(players) {
			startPlayer = 0
		}

		winningHandValue := 0
		for i, _ := range players {
			if players[i].HandValue() > winningHandValue {
				winningHandValue = players[i].HandValue()
			}
		}
		for i, _ := range players {
			if players[i].HandValue() == winningHandValue {
				players[i].roundsWon++
			}
			players[i].Forfeit()
		}
	}
	_, winner := PlayerWon()
	fmt.Println("Player won:", winner, players[winner].name)
}


func SendAll(args ...interface{}) {
	for _, player := range players {
		if !player.lost {
			player.Send(args...)
		}
	}
}

func SurvivingPlayers() int {
	result := 0
	for _, player := range players {
		if !player.lost {
			result++
		}
	}
	return result
}

func PlayerWon() (bool, int) {
	for i, player := range players {
		if player.roundsWon >= 4 {
			return true, i
		}
	}
	return false, 0
}