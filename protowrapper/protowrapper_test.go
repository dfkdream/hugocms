package protowrapper

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/dfkdream/hugocms/internal"
	"github.com/dfkdream/hugocms/user"
)

func TestEncoder_Encode(t *testing.T) {
	u := user.User{
		Id:       internal.GenerateRandomKey(8),
		Username: internal.GenerateRandomKey(8),
		Hash:     internal.GenerateRandomKey(32),
		Salt:     internal.GenerateRandomKey(32),
	}

	var writer bytes.Buffer
	err := NewEncoder(&writer).Encode(&u)
	if err != nil {
		t.Error(err)
	}

	expected, err := proto.Marshal(&u)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(expected, writer.Bytes()) {
		t.Error("Encode result not equals")
	}

}

func TestDecoder_Decode(t *testing.T) {
	u := user.User{
		Id:       internal.GenerateRandomKey(8),
		Username: internal.GenerateRandomKey(8),
		Hash:     internal.GenerateRandomKey(32),
		Salt:     internal.GenerateRandomKey(32),
	}

	buff, err := proto.Marshal(&u)
	if err != nil {
		t.Error(err)
	}

	var result user.User
	err = NewDecoder(bytes.NewReader(buff)).Decode(&result)
	if err != nil {
		t.Error(err)
	}

	if !(u.Id == result.Id && u.Username == result.Username && u.Hash == result.Hash && u.Salt == result.Salt) {
		t.Error("Decode result not equals")
	}
}
