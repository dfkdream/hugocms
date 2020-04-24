package signin

import (
	"testing"

	"github.com/dfkdream/hugocms/internal"
	"github.com/dfkdream/hugocms/user"
)

func TestHashValidatePassword(t *testing.T) {
	password := internal.GenerateRandomKey(32)
	u, err := user.New("id", "username", password, []string{"*"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !u.Validate("id", password) {
		t.Error("password validation failed")
	}
	if u.Validate("id", "0000000000000000") {
		t.Error("password validation failed")
	}
}
