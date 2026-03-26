package ui

import (
	"fmt"
	"wordGame/internal/data"
)

// Вывод приветсвенного сообщения
func PrintMsg() {
	fmt.Println("Welcome to the game")
	fmt.Printf("Guess the hidden word by the given hint\n")
	fmt.Printf("\n       >>>>>Rules<<<<<\n\n")
	fmt.Printf("1. You have 5 attempts.\n2. Find all words to win.\n3. Have fun and be smart.\n\n")
	fmt.Println("Chose the category of words to Start the Game")
	categories := data.GetCategories()

	for i, c := range categories {
		fmt.Printf("%d. %s\n", i+1, c.Name)
	}
}

func CountNumOFWords(w []data.Word) {
	fmt.Println("Great choice!\nThe number of words to guess:")
	fmt.Println(len(w))
}
