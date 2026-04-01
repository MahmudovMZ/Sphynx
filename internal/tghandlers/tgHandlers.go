package tghandlers

import (
	"fmt"
	"log"
	Game "wordGame/internal/game"
	"wordGame/internal/models"
	data "wordGame/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/bcrypt"
)

var botState = make(map[int64]string)
var gameData = make(map[int64]*Game.Game)
var userState = make(map[int64]string)
var userData = make(map[int64]map[string]string)
var bot *tgbotapi.BotAPI

func BotHandler(bot2 *tgbotapi.BotAPI, update tgbotapi.Update) {
	bot = bot2

	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	//Current bot state
	stage := botState[chatID]

	switch stage {

	// Main menu = Program Start
	case "":
		if text == "/start" {
			//Sending keyboard for choosing a option
			sendMenuKeyBoard(chatID)
			//Changing to the next state
			userState[chatID] = models.STATE_WAITING_CHOICE
		}
		switch text {
		case "Registration":
			userData[chatID] = make(map[string]string)
			send(chatID, "Speak, traveler...")
			send(chatID, "You stand before the gate of beginnings.")
			send(chatID, "To pass, you must reveal what defines you.")
			send(chatID, "First — tell me the name by which you are known.")
			send(chatID, "Choose it wisely… for it shall echo in every answer you give.")

			botState[chatID] = models.STATE_STARTING_REGISTRATION
			userState[chatID] = models.STATE_WAITING_NAME
		case "Login", "/login":
			userData[chatID] = make(map[string]string)
			send(chatID, "Welcome back, traveler... to regain your path, the gate must recognize you.")
			send(chatID, "Speak your name:")
			botState[chatID] = models.STATE_STARTING_LOGIN
			userState[chatID] = models.STATE_WAITING_LOGIN_NAME
		case "Leader Board":
			send(chatID, "leader board")

		case "Exit":
			send(chatID, "exit")

		}

	//Starting registration scenario
	case models.STATE_STARTING_REGISTRATION:
		userStage := userState[chatID]
		switch userStage {

		case models.STATE_WAITING_NAME: //getting users Name
			userData[chatID]["name"] = text
			send(chatID, fmt.Sprintf("%s... Your name is known.\nNow speak the secret that will guard your path.", userData[chatID]["name"]))
			userState[chatID] = models.STATE_WAITING_PASS

		case models.STATE_WAITING_PASS: //getting users password
			userData[chatID]["pass"] = text
			send(chatID, "Hold these secrets close... for without them, the path back will remain closed.")
			userN := userData[chatID]["name"]
			userP := userData[chatID]["pass"]
			SignUp(userN, userP)

			//deleting local data
			delete(userData, chatID)
			userState[chatID] = ""
			botState[chatID] = ""
		}

	case models.STATE_STARTING_LOGIN:
		userStage := userState[chatID]
		switch userStage {
		case models.STATE_WAITING_LOGIN_NAME:
			userData[chatID]["name"] = text
			send(chatID, fmt.Sprintf(
				"Ah, %s... I see you. Now whisper the secret that guards your way.",
				userData[chatID]["name"],
			))
			userState[chatID] = models.STATE_WAITING_LOGIN_PASS

		case models.STATE_WAITING_LOGIN_PASS:
			userData[chatID]["pass"] = text
			send(chatID, "The path is open, if your secret matches what I know... Stand by for confirmation.")

			userN := userData[chatID]["name"]
			userP := userData[chatID]["pass"]

			matchedUsers, err := Login(userN, userP)
			log.Println(matchedUsers)
			if err != nil {
				send(chatID, "An error occurred while verifying your credentials. Please try again later.")
				return
			}

			if len(matchedUsers) == 0 {
				send(chatID, "Your secret does not match. The gate remains closed.")
				userState[chatID] = models.STATE_WAITING_LOGIN_NAME
				botState[chatID] = models.STATE_STARTING_LOGIN
			} else {
				user := matchedUsers[0]
				send(chatID, fmt.Sprintf("The gate opens. You have returned successfully, %s!", user.Username))
				sendCategoryKeyboard(chatID)
				botState[chatID] = models.STATE_WAITING_CATEGORY
				userState[chatID] = ""
				delete(userData, chatID)
			}
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
		botState[chatID] = models.STATE_WAITING_GAME

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
			botState[chatID] = ""

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

func sendMenuKeyBoard(chatID int64) {
	menu := models.Menu

	//local variables in the func to prevent multiplying buttons
	rows := make([][]tgbotapi.KeyboardButton, 0)
	row := make([]tgbotapi.KeyboardButton, 0)

	for i, m := range menu {
		row = append(row, tgbotapi.NewKeyboardButton(m.Title))
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
	msg := tgbotapi.NewMessage(chatID, "Choose your path and begin.")
	msg.ReplyMarkup = keyboard

	if _, err := bot.Send(msg); err != nil {
		log.Println("Could not send category keyboard:", err)
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

func SignUp(userN, userP string) {
	// Generate bcrypt hash from plain password
	hash, err := bcrypt.GenerateFromPassword([]byte(userP), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Save the username and hashed password to the database
	err = data.SignUp_user(userN, string(hash))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nAccount with username %s has been created\n\n", userN)
}

func Login(userN, userP string) ([]models.User, error) {
	users, err := data.GetByUsersName(userN)
	if err != nil {
		return nil, fmt.Errorf("error retrieving users: %w", err)
	}

	if len(users) == 0 {
		return nil, nil
	}

	var matched []models.User
	for _, u := range users {
		err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(userP))
		if err == nil {
			matched = append(matched, u)
		}
	}

	return matched, nil
}
