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

	// Used for keeping track of the current g
	board          [BOARD_WIDTH][BOARD_HEIGHT]int
	currentHeights [BOARD_WIDTH]int
	lastMoveX      int
	lastMoveY      int

	// Store the turns in encoded format
	encodedBoard uint64
	encodedTurns []uint64
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
			turn := g.board[j][i]

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

	if turn < 0 || turn > BOARD_WIDTH {
		return false
	}

	if g.currentHeights[turn] >= BOARD_HEIGHT {
		return false
	}

	return true

}

func (g *Game) AddNewTurn(turn int) {

	// Add the turn for the player to the board
	height := BOARD_HEIGHT - 1 - g.currentHeights[turn]
	g.board[turn][height] = g.currentPlayer
	g.currentHeights[turn] += 1

	g.lastMoveX = turn
	g.lastMoveY = height

	// Switch player
	if g.currentPlayer == 1 {
		g.currentPlayer = 2
	} else {
		g.currentPlayer = 1
	}

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
}

func (g *Game) CheckForWin(turn int) bool {

	lmx := turn
	lmy := BOARD_HEIGHT - 1 - g.currentHeights[turn]

	// Check horizontal
	{
		left_to_win := 3

		stop_1 := false
		stop_2 := false
		for x := 1; x <= 4; x++ {
			if !stop_1 && lmx-x >= 0 {
				if g.board[lmx-x][lmy] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmx+x < BOARD_WIDTH {
				if g.board[lmx+x][lmy] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_2 = true
				}

			}
			if left_to_win <= 0 {
				return true
			}

		}

	}

	// Check vertical
	{
		left_to_win := 3

		stop_1 := false
		stop_2 := false

		for x := 1; x <= 4; x++ {
			if !stop_1 && lmy-x >= 0 {
				if g.board[lmx][lmy-x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy+x < BOARD_HEIGHT {
				if g.board[lmx][lmy+x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_2 = true
				}

			}
			if left_to_win <= 0 {
				return true
			}

		}

	}

	// Check diagonal bottom to top
	{
		left_to_win := 3

		stop_1 := false
		stop_2 := false

		for x := 1; x <= 4; x++ {
			if !stop_1 && lmy+x < BOARD_HEIGHT && lmx-x >= 0 {
				if g.board[lmx-x][lmy+x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy-x >= 0 && lmx+x < BOARD_WIDTH {
				if g.board[lmx+x][lmy-x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_2 = true
				}

			}
			if left_to_win <= 0 {
				return true
			}

		}

	}

	// Check diagonal top to bottom
	{
		left_to_win := 3

		stop_1 := false
		stop_2 := false

		for x := 1; x <= 4; x++ {
			if !stop_1 && lmy-x >= 0 && lmx-x >= 0 {
				if g.board[lmx-x][lmy-x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy+x < BOARD_HEIGHT && lmx+x < BOARD_WIDTH {
				if g.board[lmx+x][lmy+x] == g.currentPlayer {
					left_to_win -= 1
				} else {
					stop_2 = true
				}

			}
			if left_to_win <= 0 {
				return true
			}

		}

	}

	return false
}
