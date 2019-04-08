package users

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	// DB is the reference to our DB, which contains our user data.
	DB = newDB()

	// ErrUserAlreadyExists is the error thrown when a user attempts to create
	// a new user in the DB with a duplicate username.
	ErrUserAlreadyExists = errors.New("users: username already exists")
)

// Store is a simple in memory database, and it's protected by a read-write mutex, so no race condition
// for two different goroutines (since map is not safe for concurrency)
type Store struct {
	rwm *sync.RWMutex
	m   map[string]string
}

// newDB is for initializing in memory DB when the program starts
func newDB() *Store {
	return &Store{
		rwm: &sync.RWMutex{},
		m:   make(map[string]string),
	}
}

// NewUser accepts a username and password and creates a new user in the DB
func NewUser(username string, password string) error {
	err := exists(username)
	if err != nil {
		return err
	}

	DB.rwm.Lock()
	defer DB.rwm.Unlock()

	// will generate salted password, and it takes a byte slice, not a string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// write into the map
	DB.m[username] = string(hashedPassword)
	return nil
}

// AuthenticateUser accepts a username and password, and checks that the given password matches the hashed password
func AuthenticateUser(username string, password string) error {
	DB.rwm.RLock()
	defer DB.rwm.RUnlock()

	hashedPassword := DB.m[username]
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}

// OverrideOldPassword overrides the old password
func OverrideOldPassword(username string, password string) error {
	// Just like in NewUser
	DB.rwm.Lock()
	defer DB.rwm.Unlock()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	DB.m[username] = string(hashedPassword)
	return nil
}

// exist is an internal utility function for ensuring the usernames are unique
func exists(username string) error {
	// RLock locks rw for reading
	DB.rwm.RLock()
	defer DB.rwm.RUnlock()

	if DB.m[username] != "" {
		return ErrUserAlreadyExists
	}
	return nil
}
