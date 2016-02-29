package records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ExtendedRouterFlow struct {
	NextHopType uint32
	NextHop     net.IP `ipVersionLookUp:"NextHopType"`
	SrcMask     uint32
	DstMask     uint32
}

func (f ExtendedRouterFlow) String() string {
	type X ExtendedRouterFlow
	x := X(f)
	return fmt.Sprintf("ExtendedRouterFlow: %+v", x)
}

// RecordName returns the Name of this flow record
func (f ExtendedRouterFlow) RecordName() string {
	return "ExtendedRouterFlow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedRouterFlow) RecordType() int {
	return TypeExtendedRouterFlowRecord
}

func (f ExtendedRouterFlow) calculateBinarySize() int {
	var size int

	size += binary.Size(f.NextHopType)
	size += binary.Size(f.NextHop)
	size += binary.Size(f.SrcMask)
	size += binary.Size(f.DstMask)

	return size
}

func DecodeExtendedRouterFlow(r io.Reader) (ExtendedRouterFlow, error) {
	var err error

	f := ExtendedRouterFlow{}

	err = Decode(r, &f)

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
