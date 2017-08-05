package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("Connect 4!")

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
			break
		}
	}
}
