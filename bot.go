package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type NextTurn struct {
	next   uint64
	wins   int
	losses int
}

type Bot struct {
	indexedTurns map[uint64][]*NextTurn

	debug bool
}

func NewBot() *Bot {
	return &Bot{indexedTurns: map[uint64][]*NextTurn{}}
}

func (b *Bot) RecordGame(g *Game) {

	// Index the first turn
	b.indexTurn(0, g.encodedTurns[0], 0%2 != len(g.encodedTurns)%2)

	// Index all the other turns
	for i := 0; i < len(g.encodedTurns)-1; i++ {

		turn := g.encodedTurns[i]
		nextTurn := g.encodedTurns[i+1]

		// Did this turn end up winning the game
		won := i%2 != len(g.encodedTurns)%2

		b.indexTurn(turn, nextTurn, won)
	}
}

func (b *Bot) indexTurn(turn, nextTurn uint64, won bool) {

	wins := 0
	losses := 0

	if won {
		wins += 1
	} else {
		losses += 1
	}

	found := false
	for _, indexedNextTurn := range b.indexedTurns[turn] {
		if indexedNextTurn.next == nextTurn {

			indexedNextTurn.wins += wins
			indexedNextTurn.losses += losses

			found = true
			break
		}
	}

	if !found {
		b.indexedTurns[turn] = append(b.indexedTurns[turn], &NextTurn{next: nextTurn, wins: wins, losses: losses})
	}

}

func (b *Bot) log(message string, a ...interface{}) {
	if b.debug {
		fmt.Printf("[DEBUG] "+message, a...)
	}
}

func (b *Bot) NextMove(encodedLastMove uint64) int {

	indexedNextTurns := b.indexedTurns[encodedLastMove]

	// This decides if we should choose a random move to spread the choices over all possible choices
	spreadChoice := rand.Float32() >= 0.9

	var move int
	if spreadChoice || len(indexedNextTurns) == 0 {
		b.log("Bot randomising this turn; spreadChoice=%f lenIndexedTurns=%d\n", spreadChoice, len(indexedNextTurns))
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
			b.log("Bot found for move=%d; nextMove=%d wins=%d losses=%d winPercentage=%f\n", encodedLastMove, nextMove, nt.wins, nt.losses, winPercentage)
		}
		b.log("Bot found for move=%d; best next move; %d(%f%)\n", encodedLastMove, bestMove, bestMoveWinPercentage)
		if bestMoveWinPercentage == 0 {
			move = rand.Intn(BOARD_WIDTH)
		} else {
			// Need to decode the move to it in the range (0, BOARD_WIDTH]

			move := bestMove - encodedLastMove
			// Loop and shift until value > 0? Seems a bit shit
			for {
				if move <= 32 { // TODO(): Make this not hardcoded
					break
				}
				move >>= uint64(BOARD_HEIGHT)
			}
			move = uint64(math.Log2(float64(move)))
			if move < 0 || move >= BOARD_WIDTH {
				// TODO(): Remove this when we know the encoding / decoding is working
				fmt.Printf("pMove=%d bMove=%d move=%d; wp=%f countPossibleMoves=%d\n", encodedLastMove, bestMove, move, bestMoveWinPercentage, len(indexedNextTurns))
				panic(fmt.Sprintf("Move is a bad value=%d; bestMove=%d, previousMove=%d", move, bestMove, encodedLastMove))
			}
		}
	}

	return move
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
