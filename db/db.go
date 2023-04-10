package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.mills.io/prologic/bitcask"
)

type Db interface {
	Get(string) (*Movie, error)
	Put(string, *Movie) error
	Has(string) bool
	Close() error
}

type Rating string

type Movie struct {
	UserRating   Rating
	CriticRating Rating
}

type MovieDb struct {
	*bitcask.Bitcask
}

// OpenDB opens a Bitcask database at the specified path and returns a SimpleDB instance
func OpenDB(path string) (*MovieDb, error) {
	bc, err := bitcask.Open(path)
	if err != nil {
		return nil, err
	}
	return &MovieDb{bc}, nil
}

// Close closes the Bitcask database
func (db *MovieDb) Close() error {
	return db.Bitcask.Close()
}

func (db *MovieDb) Has(movieName string) bool {
	return db.Bitcask.Has([]byte(movieName))
}

func (db *MovieDb) Get(movieName string) (*Movie, error) {
	jsonValue, err := db.Bitcask.Get([]byte(movieName))
	if err == bitcask.ErrKeyNotFound {
		return nil, fmt.Errorf("key not found: %s", movieName)
	}
	// Unmarshal the JSON bytes to a Movie struct
	var movie Movie
	err = json.Unmarshal(jsonValue, &movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// Put associates the given value with the given key
func (db *MovieDb) Put(movieName string, movie *Movie) error {
	if movie == nil {
		return errors.New("value cannot be nil")
	}

	// Convert the value to JSON bytes
	jsonValue, err := json.Marshal(movie)
	if err != nil {
		return err
	}

	// Put the JSON bytes in the Bitcask database
	return db.Bitcask.Put([]byte(movieName), jsonValue)
}
