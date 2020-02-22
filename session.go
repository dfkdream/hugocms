package main

import (
	"sync"
	"time"
)

type session struct {
	user      *user
	ip        string
	validThru time.Time
}

type sessionDB struct {
	db                map[string]*session
	sessionExtendable bool
	tokenTTL          time.Duration
	mutex             *sync.Mutex
}

func newSessionDB(sessionExtendable bool, tokenTTL time.Duration) *sessionDB {
	return &sessionDB{
		db:                make(map[string]*session),
		sessionExtendable: sessionExtendable,
		tokenTTL:          tokenTTL,
		mutex:             new(sync.Mutex),
	}
}

func (s *sessionDB) register(user *user, ip string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := generateRandomKey(64)
	s.db[key] = &session{user, ip, time.Now().Add(s.tokenTTL)}
	return key
}

func (s *sessionDB) validate(key string, ip string) (bool, *user) {
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
