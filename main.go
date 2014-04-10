package main

import "math/rand"
import "time"
import "fmt"
import "strings"
import "flag"

var players []*Player
var debug bool

func main() {
	rand.Seed(time.Now().Unix())
	
	flag.BoolVar(&debug, "v", false, "Provide more verbose output for debugging purposes")
	flag.Parse()

	players = make([]*Player, flag.NArg())
	var startPlayer = 0
	
	var deck = new(Deck)
	
	for i, _ := range players {
			players[i] = new(Player)
	}

	for won, _ := PlayerWon(); !won; won, _ = PlayerWon() {
		deck.Init()
		for i, _ := range players {
			players[i].Init(flag.Arg(i), i, len(players), startPlayer)
			cardDrawn := deck.Draw()
			players[i].Send("draw", cardDrawn)
			players[i].AddToHand(cardDrawn)
			Debug("*** Player", i, "drew", cardDrawn, "***")
			Debug("*** Player", i, "hand:", players[i].HandString(), "***")
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
			
			Debug("*** Player", currentPlayer, "drew", cardDrawn, "***")
			Debug("*** Player", currentPlayer, "hand:", players[currentPlayer].HandString(), "***")
			
			fullMove = players[currentPlayer].Recieve()
			Debug("*** Player", currentPlayer, "move:", fullMove, "***")
			
			n, _ = fmt.Sscan(fullMove, &action, &cardStr, &target, &targetCardStr)
			
			if !strings.EqualFold(action, "play") {
				//must be "forfeit" or some invalid action
				Debug("*** ERROR: Player", currentPlayer, "made an illegal play or forfeited ***")
				Out(currentPlayer)
				goto endTurn
			}
			
			if n < 2 {
				//malformed move
				Debug("*** ERROR: Player", currentPlayer, "passed too few arguments ***")
				Out(currentPlayer)
				goto endTurn
			}

			card = ParseCard(cardStr)
			if ok := players[currentPlayer].RemoveFromHand(card); !ok {
				//trying to play a card you don't have is illegal
				Debug("*** ERROR: Player", currentPlayer, "played a card they dont have:", card, "***")
				Out(currentPlayer)
				goto endTurn
			}
			Debug("*** Player", currentPlayer, "played", card, "***")
			Debug("*** Player", currentPlayer, "hand:", players[currentPlayer].HandString(), "***")

			switch card {
			case Princess:
				SendAll("played", currentPlayer, card)
				Out(currentPlayer)
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
		fmt.Print("Round ended.\t")
		for i, p := range players {
			fmt.Print("Player ", i, " score: ", p.roundsWon, "\t")
		}
		fmt.Println()
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
	for i, p := range players {
		if p.roundsWon >= 4 {
			return true, i
		}
	}
	return false, 0
}

func Out(playerNumber int) {
	players[playerNumber].Forfeit()
	SendAll("out", playerNumber, players[playerNumber].HandString())
	Debug("*** Player", playerNumber, "out:", players[playerNumber].HandString(), "***")
}

func Debug(args ...interface{}) {
	if debug {
		fmt.Println(args...)
	}
}