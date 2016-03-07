package records

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncodeDecodeExtendedSwitchFlow(t *testing.T) {
	rec := ExtendedSwitchFlow{
		SourceVlan:          1234,
		SourcePriority:      15,
		DestinationVlan:     4321,
		DestinationPriority: 1,
	}

	b := &bytes.Buffer{}

	err := Encode(b, rec)
	if err != nil {
		t.Fatal(err)
	}

	SkipHeaderBytes(b)

	decoded, err := DecodeFlow(b, TypeExtendedSwitchFlowRecord)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(rec, decoded) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", rec, decoded)
	}
}
