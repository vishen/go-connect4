package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	BOARD_WIDTH  = 7
	BOARD_HEIGHT = 6
)

type Game struct {
	currentPlayer int

	// Used for keeping track of the current game
	board          [BOARD_WIDTH][BOARD_HEIGHT]int
	currentHeights [BOARD_WIDTH]int

	// Store the turns in encoded format
	encodedBoard uint64
	encodedTurns []uint64
}

func (g *Game) GetPlayerPretty(player int) rune {
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
				buffer.WriteRune(g.GetPlayerPretty(turn))
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
	height := BOARD_HEIGHT - 1 - g.currentHeights[turn]
	g.board[turn][height] = g.currentPlayer
	g.currentHeights[turn] += 1

	if g.currentPlayer == 1 {
		g.currentPlayer = 2
	} else {
		g.currentPlayer = 1
	}

	{
		normalizedY := BOARD_HEIGHT - height - 1
		normalizedX := BOARD_WIDTH - turn - 1

		var encodedTurn uint64 = 1
		encodedTurn <<= (uint64)((normalizedY * BOARD_WIDTH) + normalizedX)

		g.encodedBoard |= encodedTurn
		g.encodedTurns = append(g.encodedTurns, g.encodedBoard)
	}
}

func main() {
	log.Println("Connect 4!")

	game := Game{currentPlayer: 1}

	for {
		fmt.Println(game.encodedTurns)
		game.DrawBoard()

		fmt.Printf("\nEnter Move for '%s': ", string(game.GetPlayerPretty(game.currentPlayer)))
		reader := bufio.NewReader(os.Stdin)
		text, _, err := reader.ReadRune()

		if err != nil {
			fmt.Println(err)
			continue
		}

		move := (int)(text - 48)

		validMove := game.CheckIfValidTurn(move)
		if !validMove {
			fmt.Printf("> Not a valid move, try again...\n")
			continue
		}

		game.AddNewTurn(move)
	}
}
