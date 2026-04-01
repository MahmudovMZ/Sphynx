package polling

import (
	"log"
	"wordGame/internal/tghandlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartPolling(bot *tgbotapi.BotAPI) {

	log.Println("Бот запущен!")

	//Настройки на получение новых изменений в чате (сообщение тоже изменение)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	//Отлавливание (отправка запросов на получение новых апдейтов)
	updates := bot.GetUpdatesChan(updateConfig)

	//updates - канал, поэтому если пришёл новый апдейт, сюда приходит сигнал
	//И в канал сохраняется структура апдейта
	for update := range updates {
		tghandlers.BotHandler(bot, update)
	}

}
