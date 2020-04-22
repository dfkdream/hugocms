package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dfkdream/hugocms/internal"
	proto "github.com/golang/protobuf/proto"

	"github.com/boltdb/bolt"

	"golang.org/x/crypto/scrypt"
)

var (
	ErrDuplicatedUser = errors.New("duplicated user found")
)

var (
	bucketKeyUsers = []byte("users")
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

func New(id, username, password string, permissions []string) (*User, error) {
	u := User{Id: id, Username: username, Permissions: permissions}
	var err error
	u.Hash, u.Salt, err = hashPassword(password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (m User) Validate(id, password string) bool {
	return m.Id == id && validatePassword(password, m.Hash, m.Salt)
}

type DB struct {
	db *bolt.DB
}

func NewDB(db *bolt.DB) *DB {
	return &DB{
		db: db,
	}
}

func (u DB) GetAllUsers() []*User {
	result := make([]*User, 0)
	err := u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketKeyUsers)
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			usr := new(User)
			err := proto.Unmarshal(v, usr)
			if err != nil {
				return err
			}
			result = append(result, usr)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return nil
	}

	return result
}

func (u DB) GetUser(id string) *User {
	var uptr *User
	err := u.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketKeyUsers)
		if c == nil {
			return nil
		}
		uptr = new(User)
		if userData := c.Get([]byte(id)); userData != nil {
			if err := proto.Unmarshal(userData, uptr); err != nil {
				log.Println("protocol buffer unmarshal failed. falling back to json.")
				return json.Unmarshal(userData, uptr)
			}
		} else {
			uptr = nil
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
		c, err := tx.CreateBucketIfNotExists(bucketKeyUsers)
		if err != nil {
			return err
		}
		u, err := proto.Marshal(user)
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

func (u DB) DeleteUser(id string) error {
	return u.db.Update(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketKeyUsers)
		if c == nil {
			return nil
		}
		return c.Delete([]byte(id))
	})
}

func (u DB) Size() int {
	result := 0
	_ = u.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketKeyUsers)
		if b == nil {
			return nil
		}
		result = b.Stats().KeyN
		return nil
	})
	return result
}
