package main

import "fmt"
import "os"
import "strings"
import "strconv"
import "math/rand"
import "time"
import "io"

func main() {
	rand.Seed(time.Now().Unix())

	var hand = make([]Card, 0, 2)
	var str string

	var myNumber int
	str = Recieve()
	fmt.Sscanf(str, "ident %v", &myNumber)
	fmt.Println("Random Player")

	var totalPlayers int
	str = Recieve()
	fmt.Sscanf(str, "players %v", &totalPlayers)
	var players = make([]int, totalPlayers)
	for i, _ := range players {
		players[i] = i
	}
	players = append(players[:myNumber], players[myNumber:]...)

	var startPlayer int
	str = Recieve()
	fmt.Sscanf(str, "start %v", &startPlayer)

	str = Recieve()
	var card string
	fmt.Sscanf(str, "draw %v", &card)
	hand = append(hand, ParseCard(card))
	
	for {
		str = Recieve()
		var message string
		var arg, arg2 string
		fmt.Sscan(str, &message, &arg, &arg2)
		
		if strings.EqualFold(message, "player") {
			if n, _ := strconv.ParseInt(arg, 10, 0); int(n) == myNumber {
				//my turn!
				//expect a draw message next
				str = Recieve()
				fmt.Sscanf(str, "draw %v", &arg)
				hand = append(hand, ParseCard(arg))
				
				//pick a random card and play it
				i := rand.Intn(len(hand))
				c := hand[i]
				hand = append(hand[:i], hand[i+1:]...)
				switch c {
				case Princess, Minister, Priestess:
					//cards taking no arguments
					fmt.Println("play", c)
				case General, Wizard, Knight, Clown:
					//cards requiring a target
					i = rand.Intn(len(players))
					fmt.Println("play", c, players[i])
				case Soldier:
					//card requiring a target and a query
					i = rand.Intn(len(players))
					q := Card(rand.Intn(8)+1)
					fmt.Println("play", c, players[i], q)
				default:
					//something must have gone wrong
					fmt.Println("forfeit")
				}
			}
		} else if strings.EqualFold(message, "draw") {
			//add the specified card to hand
			hand = append(hand, ParseCard(arg))
		} else if strings.EqualFold(message, "swap") {
			//swap the first card in hand for the specified card 
			c := ParseCard(arg)
			hand[0] = c
		} else if strings.EqualFold(message, "discard") {
			if p, _ := strconv.ParseInt(arg, 10, 0); int(p) == myNumber {
				//remove the specified card from hand
				c := ParseCard(arg2)
				for i, v := range hand {
					if v == c {
						hand = append(hand[:i], hand[i+1:]...)
					}
				}
			}
		} else if strings.EqualFold(message, "out") {
			//can't target a player that's out so remove them from the players list
			n, _ := strconv.ParseInt(arg, 10, 0)
			for i, v := range players {
				if v == int(n) {
					players = append(players[:i], players[i+1:]...)
				}
			}
		}
		//all other messages require no handling
	}
}

func Recieve() string {
	buf := make([]byte, 1)
	str := make([]byte, 0, 20)
	n, err := os.Stdin.Read(buf)
	for ; buf[0] != '\n'; n, err = os.Stdin.Read(buf) {
		if err == nil && n > 0 {
			str = append(str, buf[0])
		} else if err != io.EOF {
			panic(err)
		}
	}
	
	return string(str)
}