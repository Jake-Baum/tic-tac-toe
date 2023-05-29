package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type parseError struct {
	message string
}

func (p parseError) Error() string {
	return p.message
}

func main() {
	b := [][]Piece{
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
		{EMPTY, EMPTY, EMPTY},
	}

	fmt.Println("Welcome to tic-tac-toe!")

	turn, numberOfTurns := X, 0
	for {
		Print(b)
		fmt.Printf("It's %s's turn.  Which square would you like to select?  ", turn)
		reader := bufio.NewReader(os.Stdin)

		text, _ := reader.ReadString('\n')

		square, err := getSquare(text)

		if err != nil {
			fmt.Printf("%s is not a valid input.\n", text)
			continue
		}

		fmt.Printf("You have selected square %d.\n", square)

		if !IsValidMove(b, square) {
			fmt.Printf("%d is not a valid square.\n", square)
			continue
		}

		b[square/3][square%3] = turn
		isWinner, winner := IsWinner(b)
		if isWinner {
			Print(b)
			fmt.Printf("Congratulations %s! You have won the game!\n", winner)
			break
		}

		if turn == X {
			turn = O
		} else {
			turn = X
		}
		numberOfTurns++

		if numberOfTurns >= 9 {
			Print(b)
			fmt.Println("It's a draw!")
			break
		}
	}
}

func getSquare(s string) (int, error) {
	if !unicode.IsDigit(rune(s[0])) {
		return -1, &parseError{"string"}
	}

	return int(s[0]) - 48, nil
}
