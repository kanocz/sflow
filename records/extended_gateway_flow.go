package records

import (
	"encoding/binary"
	"fmt"
	"net"
)

// ExtendedGatewayFlow - TypeExtendedGatewayFlowRecord
type ExtendedGatewayFlow struct {
	NextHopType          uint32
	NextHop              net.IP `xdr:"lengthField=NextHopType"`
	As                   uint32
	SrcAs                uint32
	SrcPeerAs            uint32
	DstAs                uint32 `xdr:"ignore"`
	DstPeerAs            uint32 `xdr:"ignore"`
	DstAsPathSegmentsLen uint32
	DstAsPathSegments    []ExtendedGatewayFlowASPathSegment `xdr:"lengthField=DstAsPathSegmentsLen"`
	CommunitiesLen       uint32
	Communities          []uint32 `xdr:"lengthField=CommunitiesLen"`
	LocalPref            uint32
}

// As Path Segment ordering Types
const (
	AsPathSegmentTypeUnOrdered = 1
	AsPathSegmentTypeOrdered   = 2
)

// ExtendedGatewayFlowASPathSegment defines an AS Path (either ordered or unordered)
type ExtendedGatewayFlowASPathSegment struct {
	SegType uint32 // 1: Unordered Set || 2: Ordered Set
	SegLen  uint32
	Seg     []uint32 `xdr:"lengthField=SegLen"`
}

func (f ExtendedGatewayFlow) String() string {
	type X ExtendedGatewayFlow
	x := X(f)
	return fmt.Sprintf("ExtendedGatewayFlow: %+v", x)
}

// RecordName returns the Name of this flow record
func (f ExtendedGatewayFlow) RecordName() string {
	return "ExtendedGatewayFlow"
}

// RecordType returns the ID of the sflow flow record
func (f ExtendedGatewayFlow) RecordType() int {
	return TypeExtendedGatewayFlowRecord
}

// BinarySize calculates the binary size this object will have
// This can vary as net.IP can store IPv4 or IPv6 values and the ASPath and Communities have variable length as well
func (f ExtendedGatewayFlow) BinarySize() int {
	var size int

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

// PostUnmarshal looks for an ordered AS Path to fill in DstA and DstPeerAs
func (f *ExtendedGatewayFlow) PostUnmarshal() error {
	for _, asSegment := range f.DstAsPathSegments {
		if asSegment.SegType == AsPathSegmentTypeOrdered {
			// If the AS Segment is ordered then the last Element is the DstAs and the first the DstPeerAs
			f.DstAs = asSegment.Seg[len(asSegment.Seg)-1:][0]
			f.DstPeerAs = asSegment.Seg[0:1][0]
		}
	}

	return nil
}
