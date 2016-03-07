package records

import (
	"encoding/binary"
	"fmt"
	"net"
)

// ExtendedRouterFlow - TypeExtendedRouterFlowRecord
type ExtendedRouterFlow struct {
	NextHopType uint32
	NextHop     net.IP `xdr:"lengthField=NextHopType"`
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

// BinarySize calculcated the binary size of the object since net.IP can contain IPv4 or IPv6 addresses
func (f ExtendedRouterFlow) BinarySize() int {
	var size int

	size += binary.Size(f.NextHopType)
	size += binary.Size(f.NextHop)
	size += binary.Size(f.SrcMask)
	size += binary.Size(f.DstMask)

	return size
}
