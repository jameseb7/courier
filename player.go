package main

import "os"
import "os/exec"
import "io"
import "fmt"

type Player struct{
	name string
	number int

	hand []Card
	
	ai *exec.Cmd
	toPipe io.WriteCloser
	fromPipe io.ReadCloser

	ministered bool
	protected bool
	lost bool
}

func (p *Player) Init(aiName string, playerNumber int) {
	//end any existing AI
	if p != nil {
		p.Forfeit()
	}
	
	//clear old player data
	p.ministered = false
	p.protected = false
	p.lost = false
	p.hand = make([]Card, 0, 2)
	
	//set up the AI
	if aiName == "-" {
		//the player is human
		p.toPipe = os.Stdout
		p.fromPipe = os.Stdin
		p.name = "human"
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
		
		p.name = aiName
		
		if err := p.ai.Start(); err != nil {
			panic(err)
		}
	}
	
	p.number = playerNumber
	fmt.Fprintln(p.toPipe, p.number)
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

func (p *Player) Recieve(args ...interface{}) {
	if p.lost {
		panic("Nonexistent AI")
	}
	_, err := fmt.Fscan(p.fromPipe, args...)
	if err != nil {
		panic(err)
	}
}