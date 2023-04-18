package main

import (
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"shyam/toi-scraper/db"
	"shyam/toi-scraper/notify"
	"shyam/toi-scraper/toi_api"
	"strconv"
	"sync"
)

type defaultHTTPClient struct{}

func (c *defaultHTTPClient) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type defaultGoQueryClient struct{}

func (c *defaultGoQueryClient) NewDocumentFromReader(resp io.ReadCloser) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(resp)
}

func main() {

	// configure logs
	file, err := os.OpenFile("logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	// check for env vars
	token := os.Getenv("TELEGRAM_APITOKEN")
	if token == "" {
		log.Error("Unable to find TELEGRAM_APITOKEN")
		panic(token)
	}
	chatIdStr := os.Getenv("TOI_CHAT_ID")
	if chatIdStr == "" {
		log.Error("Unable to find TOI_CHAT_ID")
		panic(chatIdStr)
	}

	// get movie reviews
	langs := []toi_api.Language{
		toi_api.Tamil,
		toi_api.Telugu,
		toi_api.Malayalam,
		toi_api.Hindi,
	}
	var allReviews []toi_api.MovieReview
	var wg sync.WaitGroup
	for _, lang := range langs {
		wg.Add(1)
		go func(lang toi_api.Language) {
			defer wg.Done()
			reviews, err := getReviews(lang)
			if err != nil {
				log.Errorf("Error getting reviews for language %v: %v", lang, err)
				return
			}
			for _, review := range reviews {
				allReviews = append(allReviews, *review)
			}
		}(lang)
	}
	wg.Wait()

	// check DB
	movieDb, err := db.OpenDB("tmp/db")
	if err != nil {
		panic(err)
	}
	var moviesToNotify []notify.Movie
	for _, review := range allReviews {
		// if movie not there, we add it to the DB
		if !movieDb.Has(review.MovieName) {
			log.Infof("Adding movie %s to the DB", review.MovieName)
			err := movieDb.Put(review.MovieName, &db.Movie{
				UserRating:   db.Rating(review.MovieUserRating),
				CriticRating: db.Rating(review.MovieCriticRating),
			})
			if err != nil {
				panic(err)
			}
			shouldNotify, err := shouldNotifyMovie(review)
			if err != nil {
				panic(err)
			}
			if shouldNotify {
				moviesToNotify = append(moviesToNotify, notify.Movie{
					MovieName:    review.MovieName,
					UserRating:   notify.Rating(review.MovieUserRating),
					CriticRating: notify.Rating(review.MovieCriticRating),
					Language:     review.Language,
				})
			}
		} else {
			log.Infof("Movie %s already exists in the DB", review.MovieName)
		}
	}

	// notify movies
	log.Infof("Notifying for %d movies", len(moviesToNotify))
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Error("Unable to find TOI API token")
		panic(err)
	}

	// -939564132
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		panic(err)
	}
	var notifier notify.Notify = notify.NewTelegramNotifier(bot, chatId)
	if moviesToNotify != nil && len(moviesToNotify) != 0 {
		err := notifier.Notify(moviesToNotify)
		if err != nil {
			panic(err)
		}
	}
}

func getReviews(lang toi_api.Language) ([]*toi_api.MovieReview, error) {
	log.Infof("Getting movie reviews for language %s", toi_api.Languages[lang])
	toiMovieApi := toi_api.NewToiMovieApi(&defaultHTTPClient{}, &defaultGoQueryClient{})
	reviews, err := toiMovieApi.GetMovieReviews(lang) // TODO: Update to get movie links/posters
	if err != nil {
		return nil, err
	}
	log.Infof("Retrieved %d movie reviews for language %s", len(reviews), toi_api.Languages[lang])
	return reviews, nil
}

func shouldNotifyMovie(review toi_api.MovieReview) (bool, error) {
	userRating, err := review.MovieUserRating.ToFloat()
	if err != nil {
		return false, err
	}
	criticRating, err := review.MovieCriticRating.ToFloat()
	if err != nil {
		return false, err
	}

	if userRating >= 3.5 || criticRating >= 3.5 {
		return true, nil
	} else {
		return false, nil
	}
}
