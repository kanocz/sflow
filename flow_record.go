package sflow

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
)

// RawPacketFlow is a raw Ethernet header flow record.
type RawPacketFlow struct {
	Protocol    uint32
	FrameLength uint32
	Stripped    uint32
	HeaderSize  uint32
	Header      []byte
	DstMAC		[6]byte
	SrcMAC		[6]byte
}

func (f RawPacketFlow) String() string {
	type X RawPacketFlow
	x := X(f)
	return fmt.Sprintf("RawPacketFlow: %+v", x)
}

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

// This Structure must match the on-wire binary protocol exactly as it is filled dynamically
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

// RecordType returns the type of flow record.
func (f RawPacketFlow) RecordType() int {
	return TypeRawPacketFlowRecord
}

func decodeRawPacketFlow(r io.Reader) (RawPacketFlow, error) {
	f := RawPacketFlow{}

	var err error

	err = binary.Read(r, binary.BigEndian, &f.Protocol)
	if err != nil {
		return f, err
	}

	err = binary.Read(r, binary.BigEndian, &f.FrameLength)
	if err != nil {
		return f, err
	}

	err = binary.Read(r, binary.BigEndian, &f.Stripped)
	if err != nil {
		return f, err
	}

	err = binary.Read(r, binary.BigEndian, &f.HeaderSize)
	if err != nil {
		return f, err
	}
	if f.HeaderSize > MaximumHeaderLength {
		return f, fmt.Errorf("sflow: header length more than %d: %d",
			MaximumHeaderLength, f.HeaderSize)
	}

	padding := (4 - f.HeaderSize) % 4
	if padding < 0 {
		padding += 4
	}

	f.Header = make([]byte, f.HeaderSize+padding)

	_, err = r.Read(f.Header)
	if err != nil {
		return f, err
	}

	// We need to consume the padded length,
	// but len(Header) should still be HeaderSize.
	f.Header = f.Header[:f.HeaderSize]

	return f, err
}

func (f RawPacketFlow) encode(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.BigEndian, uint32(f.RecordType()))
	if err != nil {
		return err
	}

	// We need to calculate encoded size of the record.
	encodedRecordLength := uint32(4 * 4) // 4 32-bit records

	// Add the length of the header padded to a multiple of 4 bytes.
	padding := (4 - f.HeaderSize) % 4
	if padding < 0 {
		padding += 4
	}

	encodedRecordLength += f.HeaderSize + padding

	err = binary.Write(w, binary.BigEndian, encodedRecordLength)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, f.Protocol)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, f.FrameLength)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, f.Stripped)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.BigEndian, f.HeaderSize)
	if err != nil {
		return err
	}

	_, err = w.Write(append(f.Header, make([]byte, padding)...))

	return err
}

// RecordType returns the type of flow record.
func (f ExtendedSwitchFlow) RecordType() int {
	return TypeExtendedSwitchFlowRecord
}

func decodedExtendedSwitchFlow(r io.Reader) (ExtendedSwitchFlow, error) {
	f := ExtendedSwitchFlow{}

	err := binary.Read(r, binary.BigEndian, &f)

	return f, err
}

func (f ExtendedSwitchFlow) encode(w io.Writer) error {
	var err error

	err = binary.Write(w, binary.BigEndian, uint32(f.RecordType()))
	if err != nil {
		return err
	}

	encodedRecordLength := uint32(4 * 4) // 4 32-bit records

	err = binary.Write(w, binary.BigEndian, encodedRecordLength)
	if err != nil {
		return err
	}

	return binary.Write(w, binary.BigEndian, f)
}

// ExtendedGatewayFlow
func (f ExtendedGatewayFlow) RecordType() int {
	return TypeExtendedGatewayFlowRecord
}

func decodeStruct(r io.Reader, s interface{}) error {
	var err error

	structure := reflect.TypeOf(s)
	data := reflect.ValueOf(s)

	if structure.Kind() == reflect.Interface || structure.Kind() == reflect.Ptr {
		structure = structure.Elem()
	}

	if data.Kind() == reflect.Interface || data.Kind() == reflect.Ptr {
		data = data.Elem()
	}

	//fmt.Printf("Decoding into %+#v\n", s)

	for i := 0; i < structure.NumField(); i++ {
		field := data.Field(i)

		//fmt.Printf("Kind: %s\n", field.Kind())

		if field.CanSet() {
			switch field.Kind() {
			case reflect.Uint32:
				var buf uint32
				if err = binary.Read(r, binary.BigEndian, &buf); err != nil {
					return err
				}
				field.SetUint(uint64(buf))
			case reflect.Slice:
				switch field.Type().Name() {
				case "IP":
					var bufferSize uint32
					NextHopType := reflect.Indirect(data).FieldByName("NextHopType").Uint()
					if NextHopType == 2 {
						bufferSize = 16
					} else {
						bufferSize = 4
					}

					buffer := make([]byte, bufferSize)
					if err = binary.Read(r, binary.BigEndian, &buffer); err != nil {
						return err
					}

					field.SetBytes(buffer)
				default:
					switch reflect.SliceOf(field.Type()).String() {
					case "[]uint32", "[][]uint32":
						key := fmt.Sprintf("%sLen", structure.Field(i).Name)
						tmp := reflect.Indirect(data).FieldByName(key)
						bufferSize := tmp.Uint()
						field.Set(reflect.MakeSlice(field.Type(), int(bufferSize), int(bufferSize)))

						// Read directly from io
						if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
							return err
						}
					default:
						key := fmt.Sprintf("%sLen", structure.Field(i).Name)
						tmp := reflect.Indirect(data).FieldByName(key)
						bufferSize := tmp.Uint()

						field.Set(reflect.MakeSlice(field.Type(), int(bufferSize), int(bufferSize)))

						for x := 0; x < int(bufferSize); x++ {
							decodeStruct(r, field.Index(x).Addr().Interface())
						}
					}
				}

			default:
				return fmt.Errorf("Unhandled Field Kind: %s", field.Kind())
			}
		}
	}

	return nil
}

func decodeExtendedGatewayFlow(r io.Reader) (ExtendedGatewayFlow, error) {
	var err error

	f := ExtendedGatewayFlow{}
	err = decodeStruct(r, &f)

	return f, err
}

func (f ExtendedGatewayFlow) encodeStruct(w io.Writer, s interface{}) error {
	var err error

	structure := reflect.TypeOf(s)
	data := reflect.ValueOf(s)

	//fmt.Printf("Encoding %+#v\n", s)

	for i := 0; i < structure.NumField(); i++ {
		field := structure.Field(i)

		switch field.Type.Kind() {
		case reflect.Uint32:
			if err = binary.Write(w, binary.BigEndian, uint32(data.FieldByIndex(field.Index).Uint())); err != nil {
				return err
			}
		case reflect.Slice:
			switch field.Type.Name() {
			case "IP":
				// We have to handle net.IP in a special way
				NextHopType := reflect.Indirect(data).FieldByName("NextHopType").Uint()
				if NextHopType == 1 {
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()[12:]); err != nil {
						return err
					}
				} else if NextHopType == 2 {
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()); err != nil {
						return err
					}
				}
			default:
				//fmt.Printf("SliceType: %s\n", reflect.SliceOf(field.Type).Elem())
				switch reflect.SliceOf(field.Type).Elem().String() {
				case "[]uint32":
					// Write directly to io
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Interface()); err != nil {
						return err
					}
				default:
					for x := 0; x < data.FieldByIndex(field.Index).Len(); x++ {
						//fmt.Printf("Slice: %+#v\n", data.FieldByIndex(field.Index).Index(x))
						f.encodeStruct(w, data.FieldByIndex(field.Index).Index(x).Interface())
					}
				}
			}
		default:
			return fmt.Errorf("Unhandled Field Kind: %s", field.Type.Kind())
		}
	}

	return err
}

func (f ExtendedGatewayFlow) encode(w io.Writer) error {
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

	err = f.encodeStruct(w, f)

	return err
}

// ExtendedRouterFlow
func (f ExtendedRouterFlow) RecordType() int {
	return TypeExtendedRouterFlowRecord
}

func decodeExtendedRouterFlow(r io.Reader) (ExtendedRouterFlow, error) {
	var err error

	f := ExtendedRouterFlow{}
	err = decodeStruct(r, &f)

	return f, err
}

func (f ExtendedRouterFlow) encodeStruct(w io.Writer, s interface{}) error {
	var err error

	structure := reflect.TypeOf(s)
	data := reflect.ValueOf(s)

	//fmt.Printf("Encoding %+#v\n", s)

	for i := 0; i < structure.NumField(); i++ {
		field := structure.Field(i)

		switch field.Type.Kind() {
		case reflect.Uint32:
			if err = binary.Write(w, binary.BigEndian, uint32(data.FieldByIndex(field.Index).Uint())); err != nil {
				return err
			}
		case reflect.Slice:
			switch field.Type.Name() {
			case "IP":
				// We have to handle net.IP in a special way
				NextHopType := reflect.Indirect(data).FieldByName("NextHopType").Uint()
				if NextHopType == 1 {
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()[12:]); err != nil {
						return err
					}
				} else if NextHopType == 2 {
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()); err != nil {
						return err
					}
				}
			default:
				//fmt.Printf("SliceType: %s\n", reflect.SliceOf(field.Type).Elem())
				switch reflect.SliceOf(field.Type).Elem().String() {
				case "[]uint32":
					// Write directly to io
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Interface()); err != nil {
						return err
					}
				default:
					for x := 0; x < data.FieldByIndex(field.Index).Len(); x++ {
						//fmt.Printf("Slice: %+#v\n", data.FieldByIndex(field.Index).Index(x))
						f.encodeStruct(w, data.FieldByIndex(field.Index).Index(x).Interface())
					}
				}
			}
		default:
			return fmt.Errorf("Unhandled Field Kind: %s", field.Type.Kind())
		}
	}

	return err
}

func (f ExtendedRouterFlow) encode(w io.Writer) error {
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

	err = f.encodeStruct(w, f)

	return err
}
