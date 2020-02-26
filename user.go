package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boltdb/bolt"

	"golang.org/x/crypto/scrypt"
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
	ID       string `json:"id"`
	Username string `json:"username"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
}

func newUser(id, username, password string) (*user, error) {
	u := user{ID: id, Username: username}
	var err error
	u.Hash, u.Salt, err = hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u user) validate(id, password string) bool {
	return u.ID == id && validatePassword(password, u.Hash, u.Salt)
}

type userDB struct {
	db *bolt.DB
}

func newUserDB(db *bolt.DB) *userDB {
	return &userDB{
		db: db,
	}
}

func (u userDB) getUser(id string) *user {
	var uptr *user
	err := u.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("users"))
		if c == nil {
			return nil
		}
		uptr = new(user)
		if userData := c.Get([]byte(id)); userData != nil {
			return json.Unmarshal(userData, uptr)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	return uptr
}

func (u *userDB) setUser(user *user) {
	err := u.db.Update(func(tx *bolt.Tx) error {
		c, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		u, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return c.Put([]byte(user.ID), u)
	})
	if err != nil {
		log.Println(err)
	}
}

func (u *userDB) addUser(user *user) error {
	if u.getUser(user.ID) == nil {
		u.setUser(user)
		return nil
	}
	return errDuplicatedUser
}

func (u userDB) size() int {
	result := 0
	_ = u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if b == nil {
			return nil
		}
		result = b.Stats().KeyN
		return nil
	})
	return result
}
