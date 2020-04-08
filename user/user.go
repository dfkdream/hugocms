package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dfkdream/hugocms/internal"

	"github.com/boltdb/bolt"

	"golang.org/x/crypto/scrypt"
)

var (
	ErrDuplicatedUser = errors.New("duplicated user found")
)

func hashPassword(password string) (string, string, error) {
	salt := internal.GenerateRandomKey(32)
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

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
}

func (u User) String() string {
	if res, err := json.Marshal(u); err == nil {
		return string(res)
	} else {
		return ""
	}
}

func New(id, username, password string) (*User, error) {
	u := User{Id: id, Username: username}
	var err error
	u.Hash, u.Salt, err = hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u User) Validate(id, password string) bool {
	return u.Id == id && validatePassword(password, u.Hash, u.Salt)
}

type DB struct {
	db *bolt.DB
}

func NewDB(db *bolt.DB) *DB {
	return &DB{
		db: db,
	}
}

func (u DB) GetUser(id string) *User {
	var uptr *User
	err := u.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("users"))
		if c == nil {
			return nil
		}
		uptr = new(User)
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

func (u DB) SetUser(user *User) {
	err := u.db.Update(func(tx *bolt.Tx) error {
		c, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		u, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return c.Put([]byte(user.Id), u)
	})
	if err != nil {
		log.Println(err)
	}
}

func (u DB) AddUser(user *User) error {
	if u.GetUser(user.Id) == nil {
		u.SetUser(user)
		return nil
	}
	return ErrDuplicatedUser
}

func (u DB) Size() int {
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
