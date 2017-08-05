package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
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

func runBotVSBotGame() {

	for i := 0; i < 10; i++ {
		g := Game{currentPlayer: 1}

		for {

			move := rand.Intn(BOARD_WIDTH)

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

		fmt.Println("######################################")
		g.DrawBoard()
		fmt.Printf("Game won by '%s' by '%s'\n", g.GetCurrentPlayer(), g.WonBy())
	}
}

func main() {
	log.Println("Connect 4!")

	runBotVSBotGame()
	// runHumanVSHumanGame()
}
