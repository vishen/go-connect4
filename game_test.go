package main

import (
	"testing"
)

func TestCheckForWinHorizontal(t *testing.T) {

	g := Game{currentPlayer: 1}

	g.board[1][1] = g.currentPlayer
	g.board[1][2] = g.currentPlayer
	g.board[1][3] = g.currentPlayer

	if !g.CheckForWin(1, 4) {
		t.Error("This should have won!")
		return
	}

	if g.wonBy != 1 {
		t.Errorf("Should have a win code of '1' not '%d (%s)'", g.wonBy, g.WonBy())
		return
	}
}

func TestCheckForWinVertical(t *testing.T) {

	g := Game{currentPlayer: 1}

	g.board[1][1] = g.currentPlayer
	g.board[2][1] = g.currentPlayer
	g.board[3][1] = g.currentPlayer

	if !g.CheckForWin(4, 1) {
		t.Error("This should have won!")
		return
	}

	if g.wonBy != 2 {
		t.Errorf("Should have a win code of '2' not '%d (%s)'", g.wonBy, g.WonBy())
		return
	}
}

func TestCheckForWinRightTopToLeftBottom(t *testing.T) {

	g := Game{currentPlayer: 1}

	g.board[1][3] = g.currentPlayer
	g.board[2][2] = g.currentPlayer
	g.board[3][1] = g.currentPlayer

	if !g.CheckForWin(4, 0) {
		t.Error("This should have won!")
		return
	}

	if g.wonBy != 4 {
		t.Errorf("Should have a win code of '4' not '%d (%s)'", g.wonBy, g.WonBy())
		return
	}
}

func TestCheckForWinLeftTopToRightBottom(t *testing.T) {

	g := Game{currentPlayer: 1}

	g.board[2][0] = g.currentPlayer
	g.board[3][1] = g.currentPlayer
	g.board[4][2] = g.currentPlayer

	if !g.CheckForWin(5, 3) {
		t.Error("This should have won!")
		return
	}

	if g.wonBy != 3 {
		t.Errorf("Should have a win code of '3' not '%d (%s)'", g.wonBy, g.WonBy())
		return
	}
}
