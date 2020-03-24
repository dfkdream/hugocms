package main

import "testing"

func TestHashValidatePassword(t *testing.T) {
	password := generateRandomKey(32)
	hash, salt, err := hashPassword(password)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !validatePassword(password, hash, salt) {
		t.Error("password validation failed")
	}
	if validatePassword(password, hash, "0000000000000000") {
		t.Error("password validation failed")
	}
}
