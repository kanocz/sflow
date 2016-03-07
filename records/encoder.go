package records

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

// Sizable Records have a custom BinarySize function because their size cannot be determined by a simple binary.Size
// They usually contain dynamic / union-like values
type Sizable interface {
	BinarySize() int
}

// Encode the given record in XDR Format and write it into w
func Encode(w io.Writer, record Record) error {
	var err error

	// Write type and length of record
	if err = binary.Write(w, binary.BigEndian, uint32(record.RecordType())); err != nil {
		return err
	}

	// Calculate Total Record Length
	var encodedRecordLength int
	if data, ok := record.(Sizable); ok {
		encodedRecordLength = data.BinarySize()
	} else {
		encodedRecordLength = binary.Size(record)
	}

	if err = binary.Write(w, binary.BigEndian, uint32(encodedRecordLength)); err != nil {
		return err
	}

	if err = encodeStruct(w, record); err != nil {
		return err
	}

	return nil
}

func encodeStruct(w io.Writer, s interface{}) error {
	var err error

	structure := reflect.TypeOf(s)
	data := reflect.ValueOf(s)

	if structure.Kind() == reflect.Interface || structure.Kind() == reflect.Ptr {
		structure = structure.Elem()
	}

	if data.Kind() == reflect.Interface || data.Kind() == reflect.Ptr {
		data = data.Elem()
	}

	//fmt.Printf("Encoding %+#v\n", s)

	for i := 0; i < structure.NumField(); i++ {
		field := structure.Field(i)
		flags := parseMarshalFlags(field)

		// Do not encode fields marked with "ignoreOnMarshal" Tags
		if flags.Ignore {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Uint8:
			if err = binary.Write(w, binary.BigEndian, uint8(data.FieldByIndex(field.Index).Uint())); err != nil {
				return err
			}
		case reflect.Uint32:
			if err = binary.Write(w, binary.BigEndian, uint32(data.FieldByIndex(field.Index).Uint())); err != nil {
				return err
			}
		case reflect.Uint64:
			if err = binary.Write(w, binary.BigEndian, uint64(data.FieldByIndex(field.Index).Uint())); err != nil {
				return err
			}
		case reflect.Slice:
			switch field.Type.Name() {
			case "IP":
				var bufferSize uint32

				ipVersion := structure.Field(i).Tag.Get("ipVersion")
				switch ipVersion {
				case "4":
					bufferSize = 4
				case "6":
					bufferSize = 16
				default:
					ipType := reflect.Indirect(data).FieldByName(flags.LengthField).Uint()
					switch ipType {
					case 1:
						bufferSize = 4
					case 2:
						bufferSize = 16
					default:
						return fmt.Errorf("Invalid Value found in ipVersionLookUp Type Field. Expected 1 or 2 and got: %d", ipType)
					}
				}

				if bufferSize == 4 && data.FieldByIndex(field.Index).Len() == 16 {
					// We write only the last 4 Bytes of the buffer (net.IP uses 16 by default even for IPv4)
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()[12:]); err != nil {
						return err
					}
				} else {
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Bytes()); err != nil {
						return err
					}
				}
			default:
				// FIXME: Add padding if necessary for opaque fields
				switch reflect.SliceOf(field.Type).Elem().String() {
				case "[]uint64", "[]uint32", "[]uint16", "[]uint8":
					// Write directly to io
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Interface()); err != nil {
						return err
					}
				default:
					for x := 0; x < data.FieldByIndex(field.Index).Len(); x++ {
						encodeStruct(w, data.FieldByIndex(field.Index).Index(x).Interface())
					}
				}
			}
		default:
			return fmt.Errorf("Unhandled Field Kind: %s", field.Type.Kind())
		}
	}

	return err
}
