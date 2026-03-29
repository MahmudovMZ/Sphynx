package main

import (
	"log"
	"os"
	"wordGame/internal/config"
	"wordGame/internal/polling"
	db "wordGame/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

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

	godotenv.Load(".env")

	var tgtoken = os.Getenv("TGBOTAPI_TOKEN")
	var startMode = os.Getenv("BOT_MODE")
	if tgtoken == "" {
		log.Fatal("Set your bot token in .env")
	}

	bot, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	switch startMode {
	case "polling":
		polling.StartPolling(bot)

	}
}
