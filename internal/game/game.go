package Game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"wordGame/internal/data"
)

// Структура описывающая состояние игры
type Game struct {
	Score        int
	Lives        int
	Words        []data.Word
	CurrentIndex int
}

// Новая игра
func NewGame(CategroyId int) *Game {
	game := Game{
		Score: 0,
		Lives: 5,
		Words: data.GetWords(CategroyId),
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
		g.Score++

		return true
	} else {
		g.Lives--
		return false
	}
}

// Проверка попыток и оставшихся слов
func (g *Game) IsGameOver() bool {
	return g.Lives <= 0 || g.CurrentIndex >= len(g.Words)
}

// Запуск всей программы
func (g *Game) Run() {
	reader := bufio.NewReader(os.Stdin)

	for !g.IsGameOver() {
		hint := g.GetCurrentHint()
		fmt.Println("Hint: ", hint)
		fmt.Println("Enter your answer")

		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.TrimSpace(userAnswer)

		if g.CheckAnswer(userAnswer) {
			fmt.Printf("Correct! Score: %v, Lives: %v\n\n", g.Score, g.Lives)
			g.CurrentIndex++
		} else {
			fmt.Printf("Wrong! Remaining attempts: %v\n\n", g.Lives)
		}
	}

	if g.Lives == 0 {
		fmt.Println("You lost!")
		fmt.Println("The correct answer was: ", g.Words[g.CurrentIndex].Answer)
	} else {
		fmt.Println("Congratulations! You guessed all words!")
	}
}
