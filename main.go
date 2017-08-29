package main

import (
	"bufio"
	"fmt"
	"log"
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

func runHumanVSBotGame(bot *Bot) {

	bot.debug = true

	humanPlayer := 1
	game := Game{currentPlayer: humanPlayer}

	encodedLastMove := uint64(0)

	fmt.Println("###################################")
	fmt.Println("Starting a new game against the bot")
	fmt.Println("###################################")
	for {

		var move int

		if game.currentPlayer != humanPlayer {
			// TODO(): Don't hardcode, but currently bot can only be player two

			for {
				move = bot.NextMove(encodedLastMove, game.GetCurrentPlayerBitmap())

				if !game.CheckIfValidTurn(move) {
					fmt.Printf("Invalid move=%d by bot...?\n", move)
					continue
				}

				break
			}
			fmt.Printf("> Bot went '%d'\n", move)
		} else {
			game.DrawBoard()
			for {

				fmt.Printf("\nEnter Move for '%s': ", game.GetCurrentPlayer())
				reader := bufio.NewReader(os.Stdin)
				text, _, err := reader.ReadRune()

				if err != nil {
					fmt.Println(err)
					continue
				}

				move = (int)(text - 48)

				validMove := game.CheckIfValidTurn(move)
				if !validMove {
					fmt.Printf("\n> Not a valid move, try again...\n\n")
					continue
				}

				break
			}
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

		encodedLastMove = game.encodedTurns[len(game.encodedTurns)-1]

	}

	bot.RecordGame(&game)
}

func trainBot(games int, bot *Bot) {

	// TODO(): I believe this can be made concurrent if we lock the `RecordGame` functionality
	fmt.Printf("Playing %d bot-vs-bot games for training...\n", games)

	for i := 0; i < games; i++ {

		g := Game{currentPlayer: 1}

		encodedLastMove := uint64(0)

		for {

			move := bot.NextMove(encodedLastMove, g.GetCurrentPlayerBitmap())

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

			encodedLastMove = g.encodedTurns[len(g.encodedTurns)-1]
		}

		bot.RecordGame(&g)
	}

	fmt.Println("Finished training bot!")

}

func main() {
	log.Println("Connect 4!")

	bot := NewBot()
	games := 100000

	trainBot(games, bot)

	for {
		runHumanVSBotGame(bot)
	}
	// runHumanVSHumanGame()
}
