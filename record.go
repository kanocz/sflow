package sflow

import (
	"errors"
	"io"
)

var (
	ErrEncodingRecord = errors.New("sflow: failed to encode record")
	ErrDecodingRecord = errors.New("sflow: failed to decode record")
)

type Record interface {
	RecordType() int
	RecordName() string
	Encode(w io.Writer) error
}
