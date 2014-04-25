Courier
=======

A game manager for the Love Letter AI competition planned by [HackSoc](https://github.com/HackSoc) at York, written in Go. Based on the specification at https://github.com/HackSoc/LoveLetterManager with some modifications introduced by [46bit](https://github.com/46bit) in [his implementation](https://github.com/46bit/postman).

Courier runs the AIs passed to it and communicates with them via stdout and stdin. AIs should read messages coming in from stdin and write messages to stdout as described in https://github.com/HackSoc/LoveLetterManager/blob/master/README.md, newlines should be sent at the end of each message as Courier requires them.

After starting, an AI should expect the following messages:
* `ident <player number>` where `<player number>` is the number of the player the AI is controlling, which is zero indexed; the AI should respond to this by printing its name to stdout (followed by a newline)
* `players <number of players>` where `<number of players>` is the total number of players in the game
* `start <starting player>` where `<starting player>` is the number of the player that goes first in the current round
* `draw <starting card>` where `<starting card>` is the initial card in the hand of the player the AI controls

At the end of a round the AI will be sent SIGKILL and a new copy of the AI will be launched at the start of the next round.

## Installation

Courier can be downloaded and installed by running, if you have the `GOPATH` environment variable set up:

```
go get github.com/jameseb7/courier
go install github.com/jameseb7/courier
```

The following command can be used instead of `go get` to download Courier as a Git repository:

```
git clone https://github.com/jameseb7/courier
```

Courier can then be found in the `$GOPATH/bin/` directory.

A random player is supplied and can be installed by running

```
go install github.com/jameseb7/courier/randomplayer
```

## Usage

Courier can be run by using a command of the form:

```
./courier [-r <ROUNDS>] [-s] [-v] <AI PATHS>
```

`<AI PATHS>` are the paths to the AIs that will be run, separated by spaces with one for each player in the game. If `-` is used as one of the AI paths then that player will be controlled by a human player with messages being sent to Courier's stdout and recieved from Courier's stdin in the same format as for AI messages.

`-r <ROUNDS>`, where `<ROUNDS>` is an integer value, sets the number of rounds that must be won to win the game, which defaults to 4

`-s`, if present, causes Courier to only run a single round of the game rather, ignoring the `-r` option

`-v`, if present, causes Courier to produce detailed output for debugging and observation purposes

An example of usage, showing random players playing a single round is:

```
./courier -v -s ./randomplayer ./randomplayer ./randomplayer ./randomplayer
```
