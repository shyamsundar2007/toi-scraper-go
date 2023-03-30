package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	log.Debug("Useful debugging information.")
	log.Info("Something noteworthy happened!")
	log.Warn("You should probably take a look at this.")
	log.Error("Something failed but I'm not quitting.")

	log.Info("Sending data to telegram")
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	/* msg := tgbotapi.NewMessage(-939564132, "Hello world")
	   if _, err := bot.Send(msg); err != nil {
	       // Note that panics are a bad way to handle errors. Telegram can
	       // have service outages or network errors, you should retry sending
	       // messages or more gracefully handle failures.
	       panic(err)
	   } */
}
