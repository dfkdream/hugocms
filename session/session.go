package session

import (
	"sync"
	"time"

	"github.com/dfkdream/hugocms/internal"

	"github.com/dfkdream/hugocms/user"
)

type session struct {
	user      *user.User
	ip        string
	validThru time.Time
}

type DB struct {
	db                map[string]*session
	sessionExtendable bool
	tokenTTL          time.Duration
	mutex             *sync.Mutex
}

func NewDB(sessionExtendable bool, tokenTTL time.Duration) *DB {
	return &DB{
		db:                make(map[string]*session),
		sessionExtendable: sessionExtendable,
		tokenTTL:          tokenTTL,
		mutex:             new(sync.Mutex),
	}
}

func (s *DB) Register(user *user.User, ip string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := internal.GenerateRandomKey(64)
	s.db[key] = &session{user, ip, time.Now().Add(s.tokenTTL)}
	return key
}

func (s *DB) Validate(key string, ip string) (bool, *user.User) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if sess := s.db[key]; sess != nil {
		if time.Now().After(sess.validThru) {
			delete(s.db, key) //Expire session
			return false, nil
		}
		if sess.ip == ip {
			if s.sessionExtendable {
				sess.validThru = time.Now().Add(s.tokenTTL) //Extend session
			}
			return true, sess.user
		}
	}
	return false, nil
}
