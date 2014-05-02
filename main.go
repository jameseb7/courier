package main

import "math"
import "math/rand"
import "math/big"
import crand "crypto/rand"
import "fmt"
import "strings"
import "flag"

var players []*Player
var debug bool
var roundsNeeded int

func main() {
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	rand.Seed(seed.Int64())

	flag.BoolVar(&debug, "v", false, "Provide more verbose output for debugging purposes")
	flag.IntVar(&roundsNeeded, "r", 4, "The number of rounds that must be won to win the game")
	single := flag.Bool("s", false, "Ignore -r and stop after running only a single round of the game.")
	flag.Parse()
	if *single {
		roundsNeeded = 1
	}

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
			var fullMove, action, cardStr, queryStr string
			var target, n int
			var card, cardDrawn Card

			if players[currentPlayer].lost {
				goto endTurn
			}

			SendAll("player", currentPlayer)
			players[currentPlayer].protected = false //protection expires on your turn

			cardDrawn = deck.Draw()
			players[currentPlayer].Send("draw", cardDrawn)
			players[currentPlayer].AddToHand(cardDrawn)
			if players[currentPlayer].HasInHand(Minister) &&
				(players[currentPlayer].HasInHand(Wizard) || players[currentPlayer].HasInHand(General)) {
				Debug("*** Player", currentPlayer, "drew", cardDrawn, "***")
				Debug("*** Player", currentPlayer, "was forced to discard the minister ***")
				_ = players[currentPlayer].RemoveFromHand(Minister)
				SendAll("played", currentPlayer, Minister)
				goto endTurn
			}

			Debug("*** Player", currentPlayer, "drew", cardDrawn, "***")
			Debug("*** Player", currentPlayer, "hand:", players[currentPlayer].HandString(), "***")

			fullMove = players[currentPlayer].Recieve()
			Debug("*** Player", currentPlayer, "move:", fullMove, "***")

			n, _ = fmt.Sscan(fullMove, &action, &cardStr, &target, &queryStr)

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

			//check the right arguments have been given
			switch card {
			case Princess, Minister, Priestess:
				if n != 2 {
					Debug("*** ERROR: Player", currentPlayer, "passed the wrong number of arguments for the played card ***")
					Out(currentPlayer)
					goto endTurn
				}
			case General, Wizard, Knight, Clown:
				if n != 3 {
					Debug("*** ERROR: Player", currentPlayer, "passed the wrong number of arguments for the played card ***")
					Out(currentPlayer)
					goto endTurn
				}

			case Soldier:
				if n != 4 {
					Debug("*** ERROR: Player", currentPlayer, "passed the wrong number of arguments for the played card ***")
					Out(currentPlayer)
					goto endTurn
				}
				if (target >= len(players)) && (target < 0) {
					Debug("*** ERROR: Player", currentPlayer, "passed a bad target argument ***")
					Out(currentPlayer)
					goto endTurn
				}
			}

			SendAll("played", currentPlayer, card)

			//check target is okay
			switch card {
			case General, Wizard, Knight, Clown, Soldier:
				if (target >= len(players)) && (target < 0) {
					Debug("*** ERROR: Player", currentPlayer, "passed a bad target argument ***")
					Out(currentPlayer)
					goto endTurn
				}
				if players[target].lost {
					Debug("*** ERROR: Player", currentPlayer, "attempted to target a player who has already lost ***")
					Out(currentPlayer)
					goto endTurn
				}
				if players[target].protected {
					Debug("*** Player", currentPlayer, "tried to target Player", target, "who was protected by a priestess ***")
					Out(currentPlayer)
					goto endTurn
				}
				if target == currentPlayer {
					Debug("*** Player", currentPlayer, "tried to target themself ***")
					var availableTargets = 0
					for _, v := range players {
						if !(v.protected || v.lost) {
							availableTargets++
						}
					}
					if availableTargets <= 1 {
						Debug("*** Player", currentPlayer, "is permitted to target themself as there are no other valid targets left ***")
					} else {
						Out(currentPlayer)
						goto endTurn
					}
				}
			}

			switch card {
			case Princess:
				Out(currentPlayer)
			case Minister:
			case General:
				players[currentPlayer].Send("swap", players[target].hand[0])
				players[target].Send("swap", players[currentPlayer].hand[0])
				players[currentPlayer].hand[0], players[target].hand[0] = players[target].hand[0], players[currentPlayer].hand[0]
				Debug("*** Player", currentPlayer, "and Player", target, "swapped hands ***")
				Debug("*** Player", currentPlayer, "hand:", players[currentPlayer].HandString(), "***")
				Debug("*** Player", target, "hand:", players[target].HandString(), "***")
			case Wizard:
				Debug("*** Player", target, "discards hand and draws a new card ***")
				SendAll("discard", target, players[target].hand[0])
				if players[target].hand[0] == Princess {
					Out(target)
					goto endTurn
				}
				if deck.IsEmpty() {
					Debug("*** Player", target, "can't draw as there are no cards left ***")
					Out(target)
					goto endTurn
				}
				players[target].hand[0] = deck.Draw()
				players[target].Send("draw", players[target].hand[0])
				Debug("*** Player", target, "hand:", players[target].HandString(), "***")
			case Priestess:
				players[currentPlayer].protected = true
			case Knight:
				Debug("*** Player", currentPlayer, "and Player", target, "compare hands ***")
				Debug("*** Player", currentPlayer, "hand:", players[currentPlayer].HandString(), "***")
				Debug("*** Player", target, "hand:", players[target].HandString(), "***")
				players[target].Send("reveal", currentPlayer, players[currentPlayer].hand[0])
				players[currentPlayer].Send("reveal", target, players[target].hand[0])
				if players[currentPlayer].HandValue() > players[target].HandValue() {
					Out(target)
				} else if players[currentPlayer].HandValue() < players[target].HandValue() {
					Out(currentPlayer)
				}
			case Clown:
				Debug("*** Player", target, "reveals a", players[target].hand[0], "from their hand to Player", currentPlayer, "***")
				players[currentPlayer].Send("reveal", target, players[target].hand[0])
			case Soldier:
				Debug("*** Player", currentPlayer, "asked Player", target, "if they have a", queryStr, "in their hand ***")
				if players[target].HasInHand(ParseCard(queryStr)) {
					Out(target)
				}
			default:
				Debug("*** ERROR: Player", currentPlayer, "played an illegal card:", cardStr, "***")
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
		if p.roundsWon >= roundsNeeded {
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
