package Game

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	ui "wordGame/internal/UI"
	"wordGame/internal/config"
	"wordGame/internal/models"
	data "wordGame/internal/repository"
	db "wordGame/pkg"
)

// Структура описывающая состояние игры
type Game struct {
	Score        int
	Lives        int
	Words        []models.Word
	CurrentIndex int
	CategoryID   int
}

// Новая игра
func NewGame(categroyId int) *Game {
	game := Game{
		Score:      0,
		Lives:      5,
		Words:      data.GetWords(categroyId),
		CategoryID: categroyId,
	}
	return &game
}

// Получить текущую подсказку
func (g *Game) GetCurrentHint() string {
	if g.CurrentIndex >= len(g.Words) {
		return ""
	}
	return g.Words[g.CurrentIndex].Hint

}

// Сравнить ответ пользователя с ответом из БД
func (g *Game) CheckAnswer(userAnswer string) bool {
	correctAnswer := g.Words[g.CurrentIndex].Answer
	if g.IsGameOver() {
		return false
	}

	if strings.TrimSpace(strings.ToLower(userAnswer)) == strings.ToLower(correctAnswer) {
		g.Score += 10

		return true
	} else {
		g.Lives--
		g.Score -= 2
		g.CurrentIndex++
		return false
	}
}

// checking the remaining attempts, score and words to guess
func (g *Game) IsGameOver() bool {
	return g.Lives == 0 || g.CurrentIndex >= len(g.Words)
}

// Запуск всей программы
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)

	for !g.IsGameOver() {
		hint := g.GetCurrentHint()
		fmt.Println("\nHint: ", hint)
		fmt.Println("Enter your answer\n")

		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)

		if g.CheckAnswer(userAnswer) {
			fmt.Printf("Correct! Score: %v, Attempts: %v\n\n", g.Score, g.Lives)
			g.CurrentIndex++
		} else {
			fmt.Printf("\nWrong! Remaining attempts: %v\n", g.Lives)
			fmt.Printf("The correct answer was: %s\n\n", g.Words[g.CurrentIndex].Answer)

		}
	}

	if g.Lives == 0 {
		fmt.Println("You lost!")
		fmt.Printf("Score: %v, Attempts: %v\n\n", g.Score, g.Lives)
	} else {
		fmt.Println("Congratulations! You guessed all words!")
	}
}

func Start() {

	err := config.ReadConfig("internal/config/config.json")
	if err != nil {
		log.Fatal("Coulds not read the config files", err)
	}

	dbData := config.GetConf().Database
	err = db.ConnectDB(dbData.Username, dbData.Password, dbData.DBName, dbData.Address)
	if err != nil {
		log.Fatal("Error occured while connecting to DB", err)
	}

	defer db.CloseDB()

	ui.PrintMsg()
	category := ui.ScanCat()
	game := NewGame(category)
	words := data.GetWords(category)
	ui.CountNumOFWords(words)
	game.Run()
}
