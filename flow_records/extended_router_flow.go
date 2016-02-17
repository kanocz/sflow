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
	//FIXME ... This is not correct and hard to calculate
	encodedRecordLength := uint32(4 + 4 + 4 + 4 + 4 + 4 + 8 + 4 + 4 + 4)

	//r := reflect.ValueOf(f)
	//fmt.Printf("Total Size: %d\n", encodedRecordLength)

	err = binary.Write(w, binary.BigEndian, uint32(encodedRecordLength))
	if err != nil {
		return err
	}

	err = Encode(w, f)

	return err
}