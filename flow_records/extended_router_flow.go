package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ExtendedRouterFlow struct {
	NextHopType uint32
	NextHop     net.IP
	SrcMask     uint32
	DstMask     uint32
}

func (f ExtendedRouterFlow) String() string {
	type X ExtendedRouterFlow
	x := X(f)
	return fmt.Sprintf("ExtendedRouterFlow: %+v", x)
}

func (f ExtendedRouterFlow) RecordType() int {
	return TypeExtendedRouterFlowRecord
}

func (f ExtendedRouterFlow) calculateBinarySize() int {
	var size int = 0

	size += binary.Size(f.NextHopType)
	size += binary.Size(f.NextHop)
	size += binary.Size(f.SrcMask)
	size += binary.Size(f.DstMask)

	return size
}

func DecodeExtendedRouterFlow(r io.Reader) (ExtendedRouterFlow, error) {
	var err error

	f := ExtendedRouterFlow{}

	flags := map[string]string{
		"ipVersionLookupField": "NextHopType",
	}
	err = Decode(r, &f, flags)

	return f, err
}

func (f ExtendedRouterFlow) Encode(w io.Writer) error {
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
