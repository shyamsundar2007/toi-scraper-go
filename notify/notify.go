package notify

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"os"
)

type Rating string

type Movie struct {
	MovieName    string
	UserRating   Rating
	CriticRating Rating
	Language     string
	Link         string
}

type Notify interface {
	Notify([]Movie) error
}

type TelegramNotifier struct {
	tgClient TgAPI
	chatId   int64
}

type TgAPI interface {
	Send(msg tgbotapi.Chattable) (tgbotapi.Message, error)
}

func NewTelegramNotifier(tgClient TgAPI, chatId int64) *TelegramNotifier {
	return &TelegramNotifier{
		tgClient: tgClient,
		chatId:   chatId,
	}
}

func (tgNotifier *TelegramNotifier) Notify(movies []Movie) error {
	// configure logs
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	if movies == nil {
		return errors.New("No movies to notify")
	}

	msgToSend := ""
	for _, movie := range movies {
		log.Debugln(movie)
		msgToSend = fmt.Sprintf("TOI Movie Review\n--------------------\n"+
			"%s (%s)\nCritic: %s\nUser: %s\nLink: %s\n--------------------\n",
			movie.MovieName, movie.Language, movie.CriticRating, movie.UserRating, movie.Link)
		msg := tgbotapi.NewMessage(tgNotifier.chatId, msgToSend)
		if _, err := tgNotifier.tgClient.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
