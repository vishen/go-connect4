package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func runHumanVSHumanGame() {

	game := Game{currentPlayer: 1}

	for {
		fmt.Println(game.encodedTurns)
		fmt.Println(game.board)
		game.DrawBoard()

		fmt.Printf("\nEnter Move for '%s': ", game.GetCurrentPlayer())
		reader := bufio.NewReader(os.Stdin)
		text, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
			continue
		}

		move := (int)(text - 48)

		validMove := game.CheckIfValidTurn(move)
		if !validMove {
			fmt.Printf("\n> Not a valid move, try again...\n\n")
			continue
		}

		if game.CompleteTurn(move) {
			game.DrawBoard()
			fmt.Printf("> Yay! '%s' won '%s' \n", game.GetCurrentPlayer(), game.WonBy())
			return
		}

		if len(game.encodedTurns) == BOARD_HEIGHT*BOARD_WIDTH {
			fmt.Printf("> Well, no one won...\n")
			return
		}

	}
}

type NextTurn struct {
	next   uint64
	wins   int
	losses int
}

func (nt *NextTurn) String() string {
	return fmt.Sprintf("[NT] %d -> %d wins, %d losses", nt.next, nt.wins, nt.losses)
}

func runBotVSBotGame() {

	rand.Seed(time.Now().UTC().UnixNano())

	debug := false
	turnsIndex := map[uint64][]*NextTurn{}

	for i := 0; i < 10000; i++ {

		g := Game{currentPlayer: 1}

		for {

			var previousMove uint64
			if len(g.encodedTurns) == 0 {
				previousMove = 0
			} else {
				previousMove = g.encodedTurns[len(g.encodedTurns)-1]
			}
			indexedNextTurns := turnsIndex[previousMove]

			// TODO(): Come up with better name
			// This decides if we should choose a random move to spread the choices over all possible choices
			randomise := rand.Float32() >= 0.9

			var move int
			if randomise || len(indexedNextTurns) == 0 {
				move = rand.Intn(BOARD_WIDTH)
			} else {
				bestMove := uint64(0)
				bestMoveWinPercentage := 0.0

				// TODO(): Rather than choosing the best move, we should weight each move then randomly
				// select a weight
				for _, nt := range indexedNextTurns {
					nextMove := nt.next
					winPercentage := float64(nt.wins) / float64(nt.wins+nt.losses)

					if winPercentage > bestMoveWinPercentage {
						bestMove = nextMove
						bestMoveWinPercentage = winPercentage
					}
				}
				if bestMoveWinPercentage == 0 {
					move = rand.Intn(BOARD_WIDTH)
				} else {
					// Need to decode the move to it in the range (0, BOARD_WIDTH]

					move := bestMove - previousMove
					// Loop and shift until value > 0? Seems a bit shit
					for {
						if move <= 32 { // TODO(): Make this not hardcoded
							break
						}
						move >>= uint64(BOARD_HEIGHT)
					}
					move = uint64(math.Log2(float64(move)))
					if move < 0 || move >= BOARD_WIDTH {
						g.DrawBoard()
						fmt.Printf("pMove=%d bMove=%d move=%d; wp=%f countPossibleMoves=%d\n", previousMove, bestMove, move, bestMoveWinPercentage, len(indexedNextTurns))
						fmt.Println(g.encodedTurns)
						panic(fmt.Sprintf("Move is a bad value=%d; bestMove=%d, previousMove=%d", move, bestMove, previousMove))
					}

				}
			}

			// TODO(): Move this into a game loop or something
			if !g.CheckIfValidTurn(move) {
				continue
			}

			if g.CompleteTurn(move) {
				break
			}

			if len(g.encodedTurns) == BOARD_HEIGHT*BOARD_WIDTH {
				break
			}

		}

		// TODO(): Very hacky for now just to get it working, refactor this into something more usable / readable
		for i := -1; i < len(g.encodedTurns)-1; i++ {

			turn := uint64(0)
			// Sooooooooo bad
			if i != -1 {
				turn = g.encodedTurns[i]
			}
			nextTurn := g.encodedTurns[i+1]

			// Did this turn end up winning the game

			// HACK
			won := i%2 != len(g.encodedTurns)%2
			if i == -1 {
				won = 0%2 != len(g.encodedTurns)%2
			}
			wins := 0
			losses := 0

			if won {
				wins += 1
			} else {
				losses += 1
			}

			found := false
			for _, indexedNextTurn := range turnsIndex[turn] {
				if indexedNextTurn.next == nextTurn {

					indexedNextTurn.wins += wins
					indexedNextTurn.losses += losses

					found = true
					break
				}
			}

			if !found {
				turnsIndex[turn] = append(turnsIndex[turn], &NextTurn{next: nextTurn, wins: wins, losses: losses})
			}

		}

	}

	if debug {
		fmt.Println()

		for turn, nextTurns := range turnsIndex {
			if turn > 1000 {
				continue
			}
			fmt.Printf("> turn=%d=\n", turn)
			for _, nt := range nextTurns {
				fmt.Printf("\t%d -> %d wins, %d losses\n", nt.next, nt.wins, nt.losses)
			}
		}
	}
}

func main() {
	log.Println("Connect 4!")

	runBotVSBotGame()
	// runHumanVSHumanGame()
}
