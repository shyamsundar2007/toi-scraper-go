package toi_api

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strconv"
)

type Rating string

type MovieReview struct {
	MovieName         string
	MovieUserRating   Rating
	MovieCriticRating Rating
	Language          string
	Link              string
}

type Language int

const (
	Tamil Language = iota
	Telugu
	Malayalam
	Hindi
)

var Languages = [...]string{
	Tamil:     "Tamil",
	Telugu:    "Telugu",
	Malayalam: "Malayalam",
	Hindi:     "Hindi",
}

var baseToiUri = "https://timesofindia.indiatimes.com"

var languageLinkMap = map[Language]string{
	Tamil:     baseToiUri + "/entertainment/tamil/movie-reviews",
	Telugu:    baseToiUri + "/entertainment/telugu/movie-reviews",
	Malayalam: baseToiUri + "/entertainment/malayalam/movie-reviews",
	Hindi:     baseToiUri + "/entertainment/hindi/movie-reviews",
}

type MovieApi interface {
	GetMovieReviews(language Language) ([]MovieReview, error)
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type GoQuery interface {
	NewDocumentFromReader(io.ReadCloser) (*goquery.Document, error)
}

type ToiMovieApi struct {
	httpClient    HTTPClient
	goQueryClient GoQuery
}

func NewToiMovieApi(httpClient HTTPClient, goQueryClient GoQuery) *ToiMovieApi {
	return &ToiMovieApi{
		httpClient:    httpClient,
		goQueryClient: goQueryClient,
	}
}

func (rating *Rating) ToFloat() (float64, error) {
	if string(*rating) == "" {
		return 0, nil
	}
	value, err := strconv.ParseFloat(string(*rating), 10)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (toiMovieApi *ToiMovieApi) GetMovieReviews(language Language) ([]*MovieReview, error) {

	url := languageLinkMap[language]
	resp, err := toiMovieApi.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := toiMovieApi.goQueryClient.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var reviews []*MovieReview
	doc.Find("#perpetualListingInitial div div.FIL_right").Each(func(i int, s *goquery.Selection) {
		reviews = extractMovieReview(s, language, reviews)
	})
	doc.Find("#perpetualListing div div.FIL_right").Each(func(i int, s *goquery.Selection) {
		reviews = extractMovieReview(s, language, reviews)
	})
	return reviews, nil
}

func extractMovieReview(s *goquery.Selection, language Language, reviews []*MovieReview) []*MovieReview {
	title := s.Find("a h3").Text()
	link, _ := s.Find("a").First().Attr("href")
	criticRating := s.Find("div div:nth-child(2) span.star_count").Text()
	userRating := s.Find("div div:nth-child(3) span.star_count").Text()
	movieReview := MovieReview{
		MovieName:         title,
		MovieUserRating:   Rating(criticRating),
		MovieCriticRating: Rating(userRating),
		Language:          Languages[language],
		Link:              baseToiUri + link,
	}
	fmt.Println("Link: " + movieReview.Link)
	reviews = append(reviews, &movieReview)
	return reviews
}

//#perpetualListingInitial > div:nth-child(2) > div.FIL_right > div > div:nth-child(2) > span.star_count
//#perpetualListing > div:nth-child(3) > div.FIL_right > div > div:nth-child(3) > span.star_count
////*[@id="perpetualListingInitial"]/div[2]/div[2]/a
