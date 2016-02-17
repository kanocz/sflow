package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
)

// ExtendedSwitchFlow is an extended switch flow record.
type ExtendedSwitchFlow struct {
	SourceVlan          uint32
	SourcePriority      uint32
	DestinationVlan     uint32
	DestinationPriority uint32
}

func (f ExtendedSwitchFlow) String() string {
	type X ExtendedSwitchFlow
	x := X(f)
	return fmt.Sprintf("ExtendedSwitchFlow: %+v", x)
}

// RecordType returns the type of flow record.
func (f ExtendedSwitchFlow) RecordType() int {
	return TypeExtendedSwitchFlowRecord
}

func (f ExtendedSwitchFlow) calculateBinarySize() int {
	var size int = 0

	size += binary.Size(f)
	/*size += binary.Size(f.NextHop)
	size += binary.Size(f.SrcMask)
	size += binary.Size(f.DstMask)*/

	return size
}

func DecodedExtendedSwitchFlow(r io.Reader) (ExtendedSwitchFlow, error) {
	f := ExtendedSwitchFlow{}

	err := binary.Read(r, binary.BigEndian, &f)

	return f, err
}

func (f ExtendedSwitchFlow) Encode(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.BigEndian, uint32(f.RecordType()))
	if err != nil {
		return err
	}

	encodedRecordLength := f.calculateBinarySize()

	err = binary.Write(w, binary.BigEndian, uint32(encodedRecordLength))
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, f)
}
