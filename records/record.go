package records

import (
	"errors"
)

// Generel Encode/Decode errros //FIXME: Are these even used?
var (
	ErrEncodingRecord = errors.New("sflow: failed to encode record")
	ErrDecodingRecord = errors.New("sflow: failed to decode record")
)

// Record defines the minimum interface for every Flow and Counter sample
type Record interface {
	RecordType() int
	RecordName() string
}
