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

	playerBitmap uint64 // Does not belong on here
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
	b.indexTurn(0, g.encodedTurns[0], g.playerOneBitmap, 0%2 != len(g.encodedTurns)%2)

	// Index all the other turns
	for i := 0; i < len(g.encodedTurns)-1; i++ {

		turn := g.encodedTurns[i]
		nextTurn := g.encodedTurns[i+1]

		var playerBitmap uint64
		if i%2 == 0 {
			playerBitmap = g.playerTwoBitmap
		} else {
			playerBitmap = g.playerOneBitmap
		}

		// Did this turn end up winning the game
		won := i%2 != len(g.encodedTurns)%2

		b.indexTurn(turn, nextTurn, playerBitmap, won)
	}
}

func (b *Bot) indexTurn(turn, nextTurn, playerBitmap uint64, won bool) {

	wins := 0
	losses := 0

	if won {
		wins += 1
	} else {
		losses += 1
	}

	found := false
	for _, indexedNextTurn := range b.indexedTurns[turn] {
		if indexedNextTurn.playerBitmap == playerBitmap || indexedNextTurn.next == nextTurn {

			indexedNextTurn.wins += wins
			indexedNextTurn.losses += losses

			found = true
			break
		}
	}

	if !found {
		b.indexedTurns[turn] = append(b.indexedTurns[turn], &NextTurn{
			next:         nextTurn,
			wins:         wins,
			losses:       losses,
			playerBitmap: playerBitmap,
		})
	}

}

func (b *Bot) log(message string, a ...interface{}) {
	if b.debug {
		fmt.Printf("[DEBUG] "+message, a...)
	}
}

func (b *Bot) NextMove(encodedLastMove, playerBitmap uint64) int {

	indexedNextTurns := b.indexedTurns[encodedLastMove]

	// This decides if we should choose a random move to spread the choices over all possible choices
	spreadChoice := rand.Float32() >= 0.9

	move := -1
	if spreadChoice || len(indexedNextTurns) == 0 {
		b.log("Bot randomising this turn; spreadChoice=%t lenIndexedTurns=%d\n", spreadChoice, len(indexedNextTurns))
		move = rand.Intn(BOARD_WIDTH)
	} else {
		bestMove := uint64(0)
		bestMoveWinPercentage := 0.0

		// TODO(): Rather than choosing the best move, we should weight each move then randomly
		// select a weight
		for _, nt := range indexedNextTurns {

			if nt.playerBitmap != playerBitmap {
				continue
			}

			nextMove := nt.next
			winPercentage := float64(nt.wins) / float64(nt.wins+nt.losses)

			if winPercentage > bestMoveWinPercentage {
				bestMove = nextMove
				bestMoveWinPercentage = winPercentage
			}
			b.log("Bot found for encodedLastMove=%d; nextMove=%d wins=%d losses=%d winPercentage=%f\n", encodedLastMove, nextMove, nt.wins, nt.losses, winPercentage)
		}
		b.log("Bot found for encodedLastmove=%d; best next move; %d(%f%%)\n", encodedLastMove, bestMove, bestMoveWinPercentage)
		if bestMoveWinPercentage == 0 {
			move = rand.Intn(BOARD_WIDTH)
		} else {
			// Need to decode the move to it in the range (0, BOARD_WIDTH]

			encodedMove := bestMove - encodedLastMove
			// b.log("decoding move: move=%d bestMove=%d encodedLastMove=%d\n", encodedMove, bestMove, encodedLastMove)
			// Loop and shift until value > 0? Seems a bit shit
			for {
				if encodedMove <= 32 { // TODO(): Make this not hardcoded
					break
				}
				encodedMove >>= uint64(BOARD_HEIGHT)
				// b.log("decoding move: move=%d\n", encodedMove)
			}
			move = BOARD_WIDTH - 1 - int(math.Log2(float64(encodedMove)))
			// b.log("Bot is going with encodedMove=%d -> move=%d\n", bestMove, move)
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
