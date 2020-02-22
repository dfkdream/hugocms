package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/crypto/scrypt"
)

var (
	errDuplicatedUser = errors.New("duplicated user found")
)

func hashPassword(password string) (string, string, error) {
	salt := generateRandomKey(32)
	hashed, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return "", "", err
	}
	return fmt.Sprintf("%x", hashed), salt, nil
}

func validatePassword(password, hash, salt string) bool {
	hashed, err := scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
	if err != nil {
		return false
	}
	return fmt.Sprintf("%x", hashed) == hash
}

type user struct {
	id       string
	username string
	hash     string
	salt     string
}

func userFromContext(ctx context.Context) (*user, bool) {
	u, ok := ctx.Value(contextKeyUser).(*user)
	return u, ok
}

func newUser(id, username, password string) (*user, error) {
	u := user{id: id, username: username}
	var err error
	u.hash, u.salt, err = hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u user) validate(id, password string) bool {
	return u.id == id && validatePassword(password, u.hash, u.salt)
}

type userDB struct {
	db    map[string]*user
	mutex *sync.Mutex
}

func newUserDB() *userDB {
	return &userDB{
		db:    make(map[string]*user),
		mutex: new(sync.Mutex),
	}
}

func (u userDB) getUser(id string) *user {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.db[id]
}

func (u *userDB) setUser(user *user) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.db[user.id] = user
}

func (u *userDB) addUser(user *user) error {
	if u.getUser(user.id) == nil {
		u.setUser(user)
		return nil
	}
	return errDuplicatedUser
}

func (u userDB) size() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return len(u.db)
}
