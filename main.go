package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
)

const (
	BOARD_HEIGHT = 6
	BOARD_WIDTH  = 7
)

type Game struct {
	board   [BOARD_WIDTH][BOARD_HEIGHT]int
	heights [BOARD_WIDTH]int
	player  string

	last_move_x int
	last_move_y int
}

func (g Game) String() string {
	var buffer bytes.Buffer

	for i := 0; i < BOARD_HEIGHT; i++ {
		for j := 0; j < BOARD_WIDTH; j++ {
			switch g.board[j][i] {
			case 0:
				buffer.WriteString("-")
			case 1:
				buffer.WriteString("X")
			case 2:
				buffer.WriteString("O")
			}
			buffer.WriteString(" ")
		}
		buffer.WriteByte('\n')
	}
	for w := 0; w < BOARD_WIDTH; w++ {
		buffer.WriteRune((rune)(48 + w))
		buffer.WriteString(" ")
	}

	buffer.WriteByte('\n')

	return buffer.String()
}

func add_to_board(turn int) bool {

	// If the row is already filled return false
	if game.heights[turn] == BOARD_HEIGHT {
		return false
	}

	height := BOARD_HEIGHT - 1 - game.heights[turn]

	var piece int
	if game.player == "X" {
		piece = 1
	} else {
		piece = 2
	}

	game.board[turn][height] = piece
	game.last_move_x = turn
	game.last_move_y = height

	game.heights[turn] += 1

	return true
}

func check_for_win() bool {
	lmx := game.last_move_x
	lmy := game.last_move_y

	piece_to_look_for := game.board[lmx][lmy]

	// Check horizontal
	{
		left_to_win := 3

		stop_1 := false
		stop_2 := false
		for x := 1; x <= 4; x++ {
			if !stop_1 && lmx-x >= 0 {
				if game.board[lmx-x][lmy] == piece_to_look_for {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmx+x < BOARD_WIDTH {
				if game.board[lmx+x][lmy] == piece_to_look_for {
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
				if game.board[lmx][lmy-x] == piece_to_look_for {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy+x < BOARD_HEIGHT {
				if game.board[lmx][lmy+x] == piece_to_look_for {
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
				if game.board[lmx-x][lmy+x] == piece_to_look_for {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy-x >= 0 && lmx+x < BOARD_WIDTH {
				if game.board[lmx+x][lmy-x] == piece_to_look_for {
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
				if game.board[lmx-x][lmy-x] == piece_to_look_for {
					left_to_win -= 1
				} else {
					stop_1 = true
				}
			}
			if !stop_2 && lmy+x < BOARD_HEIGHT && lmx+x < BOARD_WIDTH {
				if game.board[lmx+x][lmy+x] == piece_to_look_for {
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

func bot_move(bot_level int) int {

	if bot_level == 1 {
		return rand.Intn(6)
	}

	return 1
}

var (
	game Game
)

func main() {

	game = Game{}
	game.player = "X"

	var move int

	for {
		if game.player == "X" {
			// Human
			fmt.Print(game)
			fmt.Printf("\nEnter Move for '%s': ", game.player)
			reader := bufio.NewReader(os.Stdin)
			text, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
				continue
			}
			if text < '0' || text > '6' {
				fmt.Println("Invalid move, please go again.")
				continue
			}

			move = (int)(text - 48)
		} else {
			// Bot

			move = bot_move(1)
		}

		// Gameplay
		if !add_to_board(move) {
			continue
		}

		if check_for_win() {
			fmt.Print(game)
			fmt.Printf("Game over, %s won!\n", game.player)
			break
		}

		if game.player == "X" {
			game.player = "O"
		} else {
			game.player = "X"
		}

	}

}
