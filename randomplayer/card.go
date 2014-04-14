package main

import "strings"

type Card uint8

const (
	None = iota
	Soldier
	Clown
	Knight
	Priestess
	Wizard
	General
	Minister
	Princess
)

func (c Card) String() string {
	switch c {
	case Soldier:
		return "soldier"
	case Clown:
		return "clown"
	case Knight:
		return "knight"
	case Priestess:
		return "priestess"
	case Wizard:
		return "wizard"
	case General:
		return "general"
	case Minister:
		return "minister"
	case Princess:
		return "princess"
	default:
		return "<NONE>"
	}
}

func ParseCard(str string) Card {
	switch str[0] {
	case 's', 'S':
		if strings.EqualFold(str, "soldier") {
			return Soldier
		}
	case 'c', 'C':
		if strings.EqualFold(str, "clown") {
			return Clown
		}
	case 'k', 'K':
		if strings.EqualFold(str, "knight") {
			return Knight
		}
	case 'p', 'P':
		if strings.EqualFold(str, "priestess") {
			return Priestess
		}
		if strings.EqualFold(str, "princess") {
			return Princess
		}
	case 'w', 'W':
		if strings.EqualFold(str, "wizard") {
			return Wizard
		}
	case 'g', 'G':
		if strings.EqualFold(str, "general") {
			return General
		}
	case 'm', 'M':
		if strings.EqualFold(str, "minister") {
			return Minister
		}
	default:
		return None
	}
	return None
}
