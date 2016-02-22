package records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
)

// Decode an sflow packet read from 'r' into the struct given by 's' - The structs datatypes have to match the binary representation in the bytestream exactly
func Decode(r io.Reader, s interface{}) error {
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
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// We can decode these kinds directly
				field.Set(reflect.New(field.Type()).Elem())
				if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
					return err
				}
			case reflect.Array:
				// For Arrays we have to create the correct structure first but can then decode directly into them
				buf := reflect.ArrayOf(field.Len(), field.Type().Elem())
				field.Set(reflect.New(buf).Elem())
				if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
					return err
				}
			case reflect.Slice:
				// For slices we need to determine the length somehow
				switch field.Type() { // Some types (IP/HardwareAddr) are handled specifically
				case reflect.TypeOf(net.IP{}):
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
							return fmt.Errorf("Unable to determine which IP Version to read for field %s\n", structure.Field(i).Name)
						}
					}

					buffer := make([]byte, bufferSize)
					if err = binary.Read(r, binary.BigEndian, &buffer); err != nil {
						return err
					}

					field.SetBytes(buffer)
				case reflect.TypeOf(HardwareAddr{}):
					buffer := make([]byte, 6)
					if err = binary.Read(r, binary.BigEndian, &buffer); err != nil {
						return err
					}
					field.SetBytes(buffer)
				default:
					// Look up the slices length via the lengthLookUp Tag Field
					lengthField := structure.Field(i).Tag.Get("lengthLookUp")
					if lengthField == "" {
						return fmt.Errorf("Variable length slice (%s) without a defined lengthLookUp. Please specify length lookup field via struct tag: `lengthLookUp:\"fieldname\"`", structure.Field(i).Name)
					}
					bufferSize := reflect.Indirect(data).FieldByName(lengthField).Uint()

					switch field.Type().Elem().Kind() {
					case reflect.Struct, reflect.Slice, reflect.Array:
						// For slices of unspecified types we call Decode revursively for every element
						field.Set(reflect.MakeSlice(field.Type(), int(bufferSize), int(bufferSize)))

						for x := 0; x < int(bufferSize); x++ {
							Decode(r, field.Index(x).Addr().Interface())
						}
					default:
						// For slices of defined length types we can look up the length and decode directly
						field.Set(reflect.MakeSlice(field.Type(), int(bufferSize), int(bufferSize)))

						// Read directly from io
						if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
							return err
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
