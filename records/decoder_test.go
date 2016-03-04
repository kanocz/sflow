package records

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"reflect"
	"testing"
)

func SkipHeaderBytes(r io.Reader) error {
	var headerBytes [8]byte
	_, err := r.Read(headerBytes[:])
	return err
}

func TestDecodeGenericRecordStatic(t *testing.T) {
	var binaryData []byte

	testFlow := ExtendedSwitchFlow{
		SourceVlan:          1000,
		SourcePriority:      1,
		DestinationVlan:     4000,
		DestinationPriority: 10,
	}

	buffer := bytes.NewBuffer(binaryData)
	binary.Write(buffer, binary.BigEndian, &testFlow)
	resultRecord, err := Decode(buffer, TypeExtendedSwitchFlowRecord)
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}

	if !reflect.DeepEqual(testFlow, resultRecord) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", testFlow, resultRecord)
	}
}

func TestDecodeGenericRecordDynamic(t *testing.T) {
	var binaryData []byte

	testFlow := ExtendedRouterFlow{
		NextHopType: 2,
		NextHop:     net.ParseIP("2001:0db8:ac10:fe01::"), //IPv4 fails with the DeepEqual
		SrcMask:     24,
		DstMask:     24,
	}

	buffer := bytes.NewBuffer(binaryData)

	testFlow.Encode(buffer)

	SkipHeaderBytes(buffer)
	resultRecord, err := Decode(buffer, TypeExtendedRouterFlowRecord)
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}

	if !reflect.DeepEqual(testFlow, resultRecord) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", testFlow, resultRecord)
	}
}
