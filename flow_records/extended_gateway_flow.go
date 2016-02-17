package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ExtendedGatewayFlow struct {
NextHopType          uint32
NextHop              net.IP
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
	var nextHopSize uint32
	if f.NextHopType == 1 {
		//IPv4
		nextHopSize = 4
	} else {
		//IPv6
		nextHopSize = 16
	}
	//FIXME ... This is not correct and hard to calculate
	encodedRecordLength := uint32(4 + nextHopSize + 4 + 4 + 4 + 4 + 8*f.DstAsPathSegmentsLen + 4 + 4*f.CommunitiesLen + 4)

	//r := reflect.ValueOf(f)
	//fmt.Printf("Total Size: %d\n", encodedRecordLength)

	err = binary.Write(w, binary.BigEndian, uint32(encodedRecordLength))
	if err != nil {
		return err
	}

	err = Encode(w, f)

	return err
}