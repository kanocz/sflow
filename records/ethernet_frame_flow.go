package records

import (
	"fmt"
	//"bytes"
	"encoding/binary"
	"io"
)

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
	/*
	   struct ethernet_counters {
	      unsigned int dot3StatsAlignmentErrors;
	      unsigned int dot3StatsFCSErrors;
	      unsigned int dot3StatsSingleCollisionFrames;
	      unsigned int dot3StatsMultipleCollisionFrames;
	      unsigned int dot3StatsSQETestErrors;
	      unsigned int dot3StatsDeferredTransmissions;
	      unsigned int dot3StatsLateCollisions;
	      unsigned int dot3StatsExcessiveCollisions;
	      unsigned int dot3StatsInternalMacTransmitErrors;
	      unsigned int dot3StatsCarrierSenseErrors;
	      unsigned int dot3StatsFrameTooLongs;
	      unsigned int dot3StatsInternalMacReceiveErrors;
	      unsigned int dot3StatsSymbolErrors;
	   }
	*/
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

func DecodeEthernetFrameFlow(r io.Reader) (EthernetFrameFlow, error) {
	f := EthernetFrameFlow{}

	var err error

	err = binary.Read(r, binary.BigEndian, &f)
	if err != nil {
		return f, err
	}

	return f, err
}

func (f EthernetFrameFlow) Encode(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		return err
	}
	return nil
}
