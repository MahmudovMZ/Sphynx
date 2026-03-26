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
	fmt.Println("1. Capitals")
	fmt.Println("2. Vehicles")
	fmt.Println("3. Movies")
	fmt.Println("4. Animals")
	fmt.Println("5. Sports")

}

func CategoryConfirm(category int) {
	categories := map[int]string{
		1: "Capitals",
		2: "Vehicles",
		3: "Movies",
		4: "Animals",
		5: "Sports",
	}
	fmt.Printf("Great! You have choosen words from %v category\n\n", categories[category])

}

func CountNumOFWords(w []data.Word) {
	fmt.Println("Great choice!\nThe number of words to guess:")
	fmt.Println(len(w))
}
