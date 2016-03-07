package records

import (
	"fmt"
	//"bytes"
	"encoding/binary"
	"io"
)

// EthernetFrameFlow - TypeEthernetFrameFlowRecord
type EthernetFrameFlow struct {
	Dot3StatsAlignmentErrors           uint32
	Dot3StatsFCSErrors                 uint32
	Dot3StatsSingleCollisionFrames     uint32
	Dot3StatsMultipleCollisionFrames   uint32
	Dot3StatsSQETestErrors             uint32
	Dot3StatsDeferredTransmissions     uint32
	Dot3StatsLateCollisions            uint32
	Dot3StatsExcessiveCollisions       uint32
	Dot3StatsInternalMacTransmitErrors uint32
	Dot3StatsCarrierSenseErrors        uint32
	Dot3StatsFrameTooLongs             uint32
	Dot3StatsInternalMacReceiveErrors  uint32
	Dot3StatsSymbolErrors              uint32
}

func (f EthernetFrameFlow) String() string {
	type X EthernetFrameFlow
	x := X(f)
	return fmt.Sprintf("EthernetFrameFlow: %+v", x)
}

// RecordType returns the type of flow record.
func (f EthernetFrameFlow) RecordType() int {
	return TypeEthernetFrameFlowRecord
}

// RecordName returns the Name of this flow record
func (f EthernetFrameFlow) RecordName() string {
	return "EthernerFrameFlow"
}
