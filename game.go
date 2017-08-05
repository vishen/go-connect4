package main

import (
	"bytes"
	"fmt"
)

const (
	BOARD_WIDTH  = 7
	BOARD_HEIGHT = 6
)

type Game struct {
	currentPlayer int

	// Used for keeping track of the current game
	// (0,0) is top left, (5,6) is bottom right # (height, width)
	board          [BOARD_HEIGHT][BOARD_WIDTH]int
	currentHeights [BOARD_WIDTH]int
	wonBy          int // Indicates how the game was won

	// Store the turns in encoded format
	encodedBoard uint64
	encodedTurns []uint64

	debug bool
}

func (g *Game) GetCurrentPlayer() string {
	return string(g.getPlayerPretty(g.currentPlayer))
}

func (g *Game) getPlayerPretty(player int) rune {
	if player == 1 {
		return 'X'
	} else {
		return 'O'
	}
}

func (g *Game) DrawBoard() {

	var buffer bytes.Buffer

	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			turn := g.board[i][j]

			if turn == 0 {
				buffer.WriteRune('-')
			} else {
				buffer.WriteRune(g.getPlayerPretty(turn))
			}
			buffer.WriteRune(' ')
		}
		buffer.WriteRune('\n')
	}
	for w := 0; w < BOARD_WIDTH; w++ {
		buffer.WriteRune((rune)(48 + w))
		buffer.WriteRune(' ')
	}

	fmt.Println(buffer.String())
}

func (g *Game) CheckIfValidTurn(turn int) bool {

	if turn < 0 || turn >= BOARD_WIDTH {
		return false
	}

	if g.currentHeights[turn] >= BOARD_HEIGHT {
		return false
	}

	return true

}

func (g *Game) CompleteTurn(turn int) bool {
	// Returns 'true' if the turn was a winning turn...?

	// Get the current height of the turn played
	height := BOARD_HEIGHT - 1 - g.currentHeights[turn]

	{
		// Encode the current turn and save the current board state as an encoded
		// integer to the list of encodedTurns
		normalizedY := BOARD_HEIGHT - height - 1
		normalizedX := BOARD_WIDTH - turn - 1

		var encodedTurn uint64 = 1
		encodedTurn <<= (uint64)((normalizedY * BOARD_WIDTH) + normalizedX)

		g.encodedBoard |= encodedTurn
		g.encodedTurns = append(g.encodedTurns, g.encodedBoard)

	}

	// Add the turn for the player to the board
	g.board[height][turn] = g.currentPlayer

	if g.CheckForWin(height, turn) {
		return true
	}

	g.currentHeights[turn] += 1

	// Switch player
	if g.currentPlayer == 1 {
		g.currentPlayer = 2
	} else {
		g.currentPlayer = 1
	}

	return false
}

func (g *Game) WonBy() string {
	switch g.wonBy {
	case 1:
		return "horizontal"
	case 2:
		return "vertical"
	case 3:
		return "left-top-to-right-bottom diagonal"
	case 4:
		return "right-top-to-left-bottom diagonal"
	}

	return "unknown"
}

func (g *Game) CheckForWin(y, x int) bool {

	winCount := 4

	// This works, but looks kind of ugly? But can't think of simplier solution at the moment
	type getXY func(x, y, direction, currentIteration int) (int, int)
	winningWays := map[int]getXY{
		1: func(x, y, d, i int) (int, int) { // horizontal
			return x + (i * d), y
		},
		2: func(x, y, d, i int) (int, int) { // vertical
			return x, y + (i * d)
		},
		3: func(x, y, d, i int) (int, int) { // left-top-to-right-bottom diagonal
			return x + (i * d), y + (i * d)
		},
		4: func(x, y, d, i int) (int, int) { // right-top-to-left-bottom diagonal
			return x - (i * d), y + (i * d)
		},
	}

	for winCode, f := range winningWays {

		consecutive := 1

		var direction int
		for i := 0; i < 2; i++ {

			if i == 0 {
				direction = 1
			} else {
				direction = -1
			}

			for j := 1; j <= winCount; j++ {

				nx, ny := f(x, y, direction, j)

				if nx < 0 || nx >= BOARD_WIDTH || ny < 0 || ny >= BOARD_HEIGHT {
					break
				}

				if g.board[ny][nx] != g.currentPlayer {
					break
				}

				consecutive++
				if consecutive >= winCount {
					g.wonBy = winCode
					return true
				}
			}
		}
	}

	return false

}
