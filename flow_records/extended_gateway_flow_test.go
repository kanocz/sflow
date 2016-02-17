package flow_records

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestCalculateBinarySizeExtendedGatewayFlow(t *testing.T) {
	rec := ExtendedGatewayFlow{
		NextHopType:          2,
		NextHop:              net.ParseIP("2001:0db8:ac10:fe01::"), //IPv4 fails with the DeepEqual
		As:                   1234,
		SrcAs:                4321,
		SrcPeerAs:            5678,
		DstAsPathSegmentsLen: 1,
		DstAsPathSegments: []ExtendedGatewayFlowASPathSegment{{
			SegType: 1,
			SegLen:  3,
			Seg:     []uint32{1234, 4321, 65535},
		}},
		CommunitiesLen: 3,
		Communities:    []uint32{1, 18, 42011},
		LocalPref:      255,
	}

	size := rec.calculateBinarySize()
	if size != 76 {
		t.Errorf("expected\n%+#v\n, got\n%+#v", 76, size)
	}
}

func TestEncodeDecodeExtendedGatewayFlowRecord(t *testing.T) {
	rec := ExtendedGatewayFlow{
		NextHopType:          2,
		NextHop:              net.ParseIP("2001:0db8:ac10:fe01::"), //IPv4 fails with the DeepEqual
		As:                   1234,
		SrcAs:                4321,
		SrcPeerAs:            5678,
		DstAsPathSegmentsLen: 1,
		DstAsPathSegments: []ExtendedGatewayFlowASPathSegment{{
			SegType: 1,
			SegLen:  3,
			Seg:     []uint32{1234, 4321, 65535},
		}},
		CommunitiesLen: 3,
		Communities:    []uint32{1, 18, 42011},
		LocalPref:      255,
	}

	b := &bytes.Buffer{}

	err := rec.Encode(b)
	if err != nil {
		t.Fatal(err)
	}

	// Skip the header section. It's 8 bytes.
	var headerBytes [8]byte

	_, err = b.Read(headerBytes[:])
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodeExtendedGatewayFlow(b)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(rec, decoded) {
		t.Errorf("expected\n%+#v\n, got\n%+#v", rec, decoded)
	}
}
