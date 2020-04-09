package protowrapper

import (
	"io"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (enc *Encoder) Encode(v proto.Message) error {
	buff, err := proto.Marshal(v)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(buff)
	return err
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

func (dec *Decoder) Decode(v proto.Message) error {
	buff, err := ioutil.ReadAll(dec.r)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buff, v)
}
