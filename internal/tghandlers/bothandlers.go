package tghandlers

import (
	"fmt"
	"log"
	Game "wordGame/internal/game"
	"wordGame/internal/models"
	data "wordGame/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userState = make(map[int64]string)
var gameData = make(map[int64]*Game.Game)
var bot *tgbotapi.BotAPI

func Handler(bot2 *tgbotapi.BotAPI, update tgbotapi.Update) {
	bot = bot2

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	//Current user state
	stage := userState[chatID]

	switch stage {

	// Main menu = Program Start
	case "":
		if text == "/start" {
			//Sending keyboard for choosing a category
			sendCategoryKeyboard(chatID)
			//Changing to the next state
			userState[chatID] = models.STATE_WAITING_CATEGORY
		}

	//User choosing a category
	case models.STATE_WAITING_CATEGORY:
		categories := data.GetCategories()
		var catID int
		found := false

		for _, cat := range categories {
			if text == cat.Name {
				catID = cat.ID
				found = true
				break
			}
		}

		if !found {
			send(chatID, "Please choose a category using the buttons.")
			return
		}

		//Creating a game
		game := Game.NewGame(catID)
		gameData[chatID] = game
		send(chatID, fmt.Sprintf("Great choice! Number of words: %d", len(game.Words)))
		//First hint
		send(chatID, game.GetCurrentHint())
		userState[chatID] = models.STATE_WAITING_GAME

	case models.STATE_WAITING_GAME:
		//Game is going on till userState == STATE_WAITING_GAME
		game := gameData[chatID]
		correct := game.CheckAnswer(text)
		if correct {
			send(chatID, "Correct")
			game.CurrentIndex++

		} else {
			send(chatID, fmt.Sprintf("Wrong answer! Remaining attempts : %d", game.Lives))
		}
		if game.IsGameOver() {
			send(chatID, fmt.Sprintf("Game over! Score: %d. Attempts: %d", game.Score, game.Lives))
			delete(gameData, chatID)
			userState[chatID] = ""

			return
		}
		//Endless cycle till quizes are end
		send(chatID, game.GetCurrentHint())

	}

}

func send(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Could not send an answer", err)
		return
	}
}

func sendCategoryKeyboard(chatID int64) {
	categories := data.GetCategories()

	//local variables in the func to prevent multiplying buttons
	rows := make([][]tgbotapi.KeyboardButton, 0)
	row := make([]tgbotapi.KeyboardButton, 0)

	for i, cat := range categories {
		row = append(row, tgbotapi.NewKeyboardButton(cat.Name))
		if (i+1)%2 == 0 {
			//adding row copy into the rows
			rows = append(rows, append([]tgbotapi.KeyboardButton{}, row...))
			row = []tgbotapi.KeyboardButton{}
		}
	}

	if len(row) > 0 {
		rows = append(rows, append([]tgbotapi.KeyboardButton{}, row...))
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	msg := tgbotapi.NewMessage(chatID, "Choose a category:")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println("Could not send category keyboard:", err)
	}
}
