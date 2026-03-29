package tghandlers

import (
	"fmt"
	"log"
	"strconv"
	Game "wordGame/internal/game"
	"wordGame/internal/models"
	data "wordGame/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//Start = вывод категорий -> ввод от пользователя = Кнопки с выбором

// var userState = make(map[int64]string)
// var userData = make(map[int64]map[string]string)
// var bot *tgbotapi.BotAPI

// func Handler(bot2 *tgbotapi.BotAPI, update tgbotapi.Update) {
// 	bot = bot2

// 	if update.Message == nil {
// 		return
// 	}

// 	chatID := update.Message.Chat.ID
// 	text := update.Message.Text

// 	//Пустое место для функции которая обрабатывает текст сообщения
// 	stage := userState[chatID]

// 	switch stage {
// 	//Главная страница (человек никуда не нажал)
// 	case "":
// 		//Проверка что человек нажал на главном экране
// 		switch text {
// 		case "/start":
// 			userData[chatID] = make(map[string]string)
// 			send(chatID, "Привет! Чё нада?")

// 		case "/registration":
// 			send(chatID, "Введите номер телефона:")
// 			userState[chatID] = models.STATE_WAITING_PHONE
// 			log.Println("registration", userState)

// 		case "/me":
// 			name := userData[chatID]["name"]
// 			phone := userData[chatID]["phone"]
// 			send(chatID, fmt.Sprintf("Ваша информация\nИмя: %s\nНомер телефона:%s", name, phone))
// 		}
// 	case models.STATE_WAITING_PHONE:

// 		if update.Message.Contact == nil {

// 			//Содержимое текста
// 			if strings.Contains(text, "+") {
// 				userState[chatID] = models.STATE_WAITING_NAME
// 				send(chatID, "Контакт сохранён. Введите Ваше имя:")
// 				userData[chatID]["phone"] = text
// 				return
// 			}

// 			if !checkPhoneNumber(text) {
// 				send(chatID, "Номер ввендён некорректно")
// 				return
// 			}

// 			userData[chatID]["phone"] = text
// 		} else {
// 			send(chatID, "Контакт сохранён. Введите Ваше имя:")
// 			userState[chatID] = models.STATE_WAITING_NAME
// 			userData[chatID]["phone"] = update.Message.Contact.PhoneNumber
// 			return
// 		}

// 		sendKeyboard(chatID, "Выберите Вашу страну", models.CountryKeyboard)
// 		userState[chatID] = models.STATE_WAITING_COUNTRY
// 	case models.STATE_WAITING_COUNTRY:

// 		switch text {
// 		case models.BUTTON_COUNTRY_TJ:
// 			userData[chatID]["phone"] = "+992" + userData[chatID]["phone"]
// 		case models.BUTTON_COUNTRY_UZ:
// 			userData[chatID]["phone"] = "+998" + userData[chatID]["phone"]
// 		}

// 		userState[chatID] = models.STATE_WAITING_NAME
// 		sendKeyboard(chatID, "Сохранено! Введите Ваше имя:", tgbotapi.NewRemoveKeyboard(true))

// 	case models.STATE_WAITING_NAME:
// 		userData[chatID]["name"] = text
// 		send(chatID, "Ваши данные сохранены!")
// 		delete(userState, chatID)
// 	}

// }

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

	// Текущее состояние пользователя
	stage := userState[chatID]

	switch stage {

	// Пользователь только начал, главное меню
	case "":
		if text == "/start" {
			categories := data.GetCategories()
			var str string
			for _, v := range categories {
				str += fmt.Sprintf("%d. %s\n", v.ID, v.Name)
			}
			// userData[chatID] = make(map[string]string)
			send(chatID, "Choose the category using numbers")
			send(chatID, str)

			// Переходим к следующему состоянию
			userState[chatID] = models.STATE_WAITING_CATEGORY
		}

	// Пользователь выбирает категорию
	case models.STATE_WAITING_CATEGORY:
		// Проверяем, что пользователь ввёл число
		catId, err := strconv.Atoi(text)
		if err != nil {
			send(chatID, "Please enter a valid number for the category.")
			return
		}

		// Проверка диапазона (например, категории с 1 по N)
		categories := data.GetCategories()
		if catId < 1 || catId > len(categories) {
			send(chatID, "Please choose a valid category number from the list.")
			return
		}

		// Создаём новую игру
		game := Game.NewGame(catId)
		gameData[chatID] = game
		// Формируем сообщение с количеством слов
		response := fmt.Sprintf("Great choice! The number of words to guess: %d", len(game.Words))
		send(chatID, response)
		send(chatID, game.GetCurrentHint())
		userState[chatID] = models.STATE_WAITING_GAME

	case models.STATE_WAITING_GAME:
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

// func sendKeyboard(chatID int64, message string, keyboard interface{}) {
// 	msg := tgbotapi.NewMessage(chatID, message)
// 	msg.ReplyMarkup = keyboard
// 	_, err := bot.Send(msg)
// 	if err != nil {
// 		log.Println("Could not send keybord", err)
// 		return
// 	}

// }

func buildReply(text string) string {
	switch text {
	case "/start":
		return "Привет! Я не знаю пока много команд, но скоро буду расти!"
	case "/help":
		return "Список доступных команд: /help, /start"
	case "/lol":
		return "Just for fun"
	}
	return "Я не знаю такой команды"
}
