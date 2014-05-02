package main

import "os"
import "os/exec"
import "io"
import "fmt"

type Player struct {
	name   string
	number int

	hand []Card

	ai       *exec.Cmd
	toPipe   io.WriteCloser
	fromPipe io.ReadCloser

	protected bool
	lost      bool

	roundsWon int
}

func (p *Player) Init(aiName string, playerNumber, nPlayers, startPlayer int) {
	//end any existing AI
	if p != nil {
		p.Forfeit()
	}

	//clear old player data
	p.protected = false
	p.lost = false
	p.hand = make([]Card, 0, 2)

	//set up the AI
	if aiName == "-" {
		//the player is human
		p.toPipe = os.Stdout
		p.fromPipe = os.Stdin
	} else {
		p.ai = exec.Command(aiName)

		if pipe, err := p.ai.StdinPipe(); err != nil {
			panic(err)
		} else {
			p.toPipe = pipe
		}

		if pipe, err := p.ai.StdoutPipe(); err != nil {
			panic(err)
		} else {
			p.fromPipe = pipe
		}

		if err := p.ai.Start(); err != nil {
			panic(err)
		}
	}

	p.number = playerNumber
	p.Send("ident", p.number)
	p.name = p.Recieve()
	p.Send("players", nPlayers)
	p.Send("start", startPlayer)
}

func (p *Player) Forfeit() {
	p.lost = true

	if p.ai != nil {
		if p.ai.Process != nil {
			p.ai.Process.Kill()
		}
		p.ai = nil
	}
	p.toPipe = nil
	p.fromPipe = nil
}

func (p *Player) Send(args ...interface{}) {
	if p.lost {
		panic("Nonexistent AI")
	}
	_, err := fmt.Fprintln(p.toPipe, args...)
	if err != nil {
		panic(err)
	}
}

func (p *Player) Recieve() string {
	if p.lost {
		panic("Nonexistent AI")
	}

	buf := make([]byte, 1)
	str := make([]byte, 0, 20)
	n, err := p.fromPipe.Read(buf)
	for ; buf[0] != '\n'; n, err = p.fromPipe.Read(buf) {
		if err == nil && n > 0 {
			str = append(str, buf[0])
		} else if err != io.EOF {
			panic(err)
		}
	}

	return string(str)
}

func (p *Player) AddToHand(c Card) {
	p.hand = append(p.hand, c)
}

func (p *Player) RemoveFromHand(c Card) (ok bool) {
	for i, v := range p.hand {
		if v == c {
			p.hand = append(p.hand[:i], p.hand[i+1:]...)
			return true
		}
	}
	return false
}

func (p *Player) HasInHand(c Card) bool {
	for _, v := range p.hand {
		if v == c {
			return true
		}
	}
	return false
}

func (p *Player) HandString() string {
	str := make([]byte, 0, 20)
	for i, v := range p.hand {
		if i != 0 {
			str = append(str, ' ')
		}
		str = append(str, []byte(v.String())...)
	}
	return string(str)
}

func (p *Player) HandValue() int {
	if p.lost {
		return 0
	}
	value := 0
	for _, c := range p.hand {
		value += int(c)
	}
	return value
}
