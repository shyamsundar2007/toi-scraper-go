package db

import (
	"fmt"
	"git.mills.io/prologic/bitcask"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMovieDb(t *testing.T) {
	// Create a Bitcask database in a temporary location
	bc, err := bitcask.Open("tmp/db")
	assert.NoError(t, err)
	defer bc.Close()
	defer func() {
		err := os.RemoveAll("tmp")
		if err != nil {
			fmt.Println("Error removing temp DB")
		}
	}()

	// Create a movie DB instance
	movieDb := MovieDb{Bitcask: bc}

	// Create test movie object
	movie := Movie{
		UserRating:   Rating("3.5"),
		CriticRating: Rating("4.0"),
	}

	// Test put method
	err = movieDb.Put("testMovie", &movie)
	assert.NoError(t, err)

	// Test Get method
	got, err := movieDb.Get("testMovie")
	assert.NoError(t, err)
	assert.Equal(t, &movie, got)

	// Test Has method
	has := movieDb.Has("testMovie")
	assert.True(t, has)

	has = movieDb.Has("movieNotExists")
	assert.False(t, has)

	// Test Close method
	err = movieDb.Close()
	assert.NoError(t, err)
}
