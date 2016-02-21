package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Encode(w io.Writer, s interface{}) error {
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
				var bufferSize uint32

				ipVersion := structure.Field(i).Tag.Get("ipVersion")
				switch ipVersion {
				case "4":
					bufferSize = 4
				case "6":
					bufferSize = 16
				default:
					lookupField := structure.Field(i).Tag.Get("ipVersionLookUp")
					switch lookupField {
					default:
						ipType := reflect.Indirect(data).FieldByName(lookupField).Uint()
						switch ipType {
						case 1:
							bufferSize = 4
						case 2:
							bufferSize = 16
						default:
							return fmt.Errorf("Invalid Value found in ipVersionLookUp Type Field. Expected 1 or 2 and got: %d", ipType)
						}
					case "":
						return fmt.Errorf("Unable to determine which IP Version to read for field %s\n", field.Type.Name())
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
				switch reflect.SliceOf(field.Type).Elem().String() {
				case "[]uint32":
					// Write directly to io
					if err = binary.Write(w, binary.BigEndian, data.FieldByIndex(field.Index).Interface()); err != nil {
						return err
					}
				default:
					for x := 0; x < data.FieldByIndex(field.Index).Len(); x++ {
						Encode(w, data.FieldByIndex(field.Index).Index(x).Interface())
					}
				}
			}
		default:
			return fmt.Errorf("Unhandled Field Kind: %s", field.Type.Kind())
		}
	}

	return err
}
