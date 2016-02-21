package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ExtendedGatewayFlow struct {
	NextHopType          uint32
	NextHop              net.IP `ipVersionLookUp:"NextHopType"`
	As                   uint32
	SrcAs                uint32
	SrcPeerAs            uint32
	DstAsPathSegmentsLen uint32
	DstAsPathSegments    []ExtendedGatewayFlowASPathSegment
	CommunitiesLen       uint32
	Communities          []uint32
	LocalPref            uint32
}

type ExtendedGatewayFlowASPathSegment struct {
	SegType uint32 // 1: Unordered Set || 2: Ordered Set
	SegLen  uint32
	Seg     []uint32
}

func (f ExtendedGatewayFlow) String() string {
	type X ExtendedGatewayFlow
	x := X(f)
	return fmt.Sprintf("ExtendedGatewayFlow: %+v", x)
}

// ExtendedGatewayFlow
func (f ExtendedGatewayFlow) RecordType() int {
	return TypeExtendedGatewayFlowRecord
}

func (f ExtendedGatewayFlow) calculateBinarySize() int {
	var size int = 0

	size += binary.Size(f.NextHopType)
	size += binary.Size(f.NextHop)
	size += binary.Size(f.As)
	size += binary.Size(f.SrcAs)
	size += binary.Size(f.SrcPeerAs)
	size += binary.Size(f.DstAsPathSegmentsLen)
	for _, segment := range f.DstAsPathSegments {
		size += binary.Size(segment.SegType)
		size += binary.Size(segment.SegLen)
		size += binary.Size(segment.Seg)
	}
	size += binary.Size(f.CommunitiesLen)
	size += binary.Size(f.Communities)
	size += binary.Size(f.LocalPref)

	return size
}

func DecodeExtendedGatewayFlow(r io.Reader) (ExtendedGatewayFlow, error) {
	var err error

	f := ExtendedGatewayFlow{}

	err = Decode(r, &f)

	return f, err
}

func (f ExtendedGatewayFlow) Encode(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.BigEndian, uint32(f.RecordType()))
	if err != nil {
		fmt.Printf("error: %s", err)
		return err
	}

	// Calculate Total Record Length
	encodedRecordLength := f.calculateBinarySize()

	err = binary.Write(w, binary.BigEndian, uint32(encodedRecordLength))
	if err != nil {
		return err
	}

	err = Encode(w, f)

	return err
}
