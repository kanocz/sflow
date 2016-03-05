package records

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"reflect"
)

type PostDecoder interface {
	PostDecode() error
}

func DecodeFlow(r io.Reader, recordType uint32) (Record, error) {
	var err error

	switch recordType {
	case TypeRawPacketFlowRecord:
		return DecodeRawPacketFlow(r)
	default:
		if recordStruct, found := flowRecordTypes[recordType]; found {
			data := reflect.New(reflect.TypeOf(recordStruct)).Elem()

			_, err = decodeInto(r, data.Addr().Interface())

			// Some records calculate extra data from the decoded values
			if data, ok := data.Addr().Interface().(PostDecoder); ok {
				data.PostDecode()
			}

			return data.Interface().(Record), err
		}
	}

	return nil, fmt.Errorf("Flow record type %d is not implemented yet\n", recordType)
}

func DecodeCounter(r io.Reader, recordType uint32) (Record, error) {
	var err error

	switch recordType {
	default:
		if recordStruct, found := counterRecordTypes[recordType]; found {
			data := reflect.New(reflect.TypeOf(recordStruct)).Elem()

			_, err = decodeInto(r, data.Addr().Interface())

			// Some records calculate extra data from the decoded values
			if data, ok := data.Addr().Interface().(PostDecoder); ok {
				data.PostDecode()
			}

			return data.Interface().(Record), err
		}
	}

	return nil, fmt.Errorf("Counter record type %d is not implemented yet\n", recordType)
}

// Decode an sflow packet read from 'r' into the struct given by 's' - The structs datatypes have to match the binary representation in the bytestream exactly
func decodeInto(r io.Reader, s interface{}) (int, error) {
	var err error
	var bytesRead int

	// If the provided datastructure has a static size we can decode it directly
	if size := binary.Size(s); size != -1 {
		err = binary.Read(r, binary.BigEndian, s)
		return size, err
	}

	structure := reflect.TypeOf(s)
	data := reflect.ValueOf(s)

	if structure.Kind() == reflect.Interface || structure.Kind() == reflect.Ptr {
		structure = structure.Elem()
	}

	if data.Kind() == reflect.Interface || data.Kind() == reflect.Ptr {
		data = data.Elem()
	}

	//fmt.Printf("Decoding into %T - %+#v\n", s, s)

	for i := 0; i < structure.NumField(); i++ {
		field := data.Field(i)

		// Do not decode fields marked with "ignoreOnMarshal" Tags
		if ignoreField := structure.Field(i).Tag.Get("ignoreOnMarshal"); ignoreField == "true" {
			continue
		}

		//fmt.Printf("Kind: %s - %s\n", field.Kind(), field.CanSet())
		//fmt.Printf("State: %s\n", s)

		if field.CanSet() {
			switch field.Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// We can decode these kinds directly
				field.Set(reflect.New(field.Type()).Elem())
				if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
					return bytesRead, err
				}
				bytesRead += binary.Size(field.Addr().Interface())
			case reflect.Array:
				// For Arrays we have to create the correct structure first but can then decode directly into them
				buf := reflect.ArrayOf(field.Len(), field.Type().Elem())
				field.Set(reflect.New(buf).Elem())
				if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
					return bytesRead, err
				}
				bytesRead += binary.Size(field.Addr().Interface())
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
								return bytesRead, fmt.Errorf("Invalid Value found in ipVersionLookUp Type Field. Expected 1 or 2 and got: %d", ipType)
							}
						case "":
							return bytesRead, fmt.Errorf("Unable to determine which IP Version to read for field %s\n", structure.Field(i).Name)
						}
					}

					buffer := make([]byte, bufferSize)
					if err = binary.Read(r, binary.BigEndian, &buffer); err != nil {
						return bytesRead, err
					}
					bytesRead += binary.Size(&buffer)

					field.SetBytes(buffer)
				case reflect.TypeOf(HardwareAddr{}):
					buffer := make([]byte, 6)
					if err = binary.Read(r, binary.BigEndian, &buffer); err != nil {
						return bytesRead, err
					}
					bytesRead += binary.Size(&buffer)
					field.SetBytes(buffer)
				default:
					// Look up the slices length via the lengthLookUp Tag Field
					lengthField := structure.Field(i).Tag.Get("lengthLookUp")
					if lengthField == "" {
						return bytesRead, fmt.Errorf("Variable length slice (%s) without a defined lengthLookUp. Please specify length lookup field via struct tag: `lengthLookUp:\"fieldname\"`", structure.Field(i).Name)
					}
					bufferSize := reflect.Indirect(data).FieldByName(lengthField).Uint()

					if bufferSize > 0 {
						switch field.Type().Elem().Kind() {
						case reflect.Struct, reflect.Slice, reflect.Array:
							// For slices of unspecified types we call Decode revursively for every element
							field.Set(reflect.MakeSlice(field.Type(), int(bufferSize), int(bufferSize)))

							for x := 0; x < int(bufferSize); x++ {
								decodeInto(r, field.Index(x).Addr().Interface())
							}
						default:
							//Apply padding
							size := bufferSize + (4-(bufferSize%4))%4

							// For slices of defined length types we can look up the length and decode directly
							field.Set(reflect.MakeSlice(field.Type(), int(size), int(size)))

							// Read directly from io
							if err = binary.Read(r, binary.BigEndian, field.Addr().Interface()); err != nil {
								return bytesRead, err
							}
							bytesRead += binary.Size(field.Addr().Interface())
						}
					}
				}
			case reflect.Struct:
				// For structs we call Decode revursively
				field.Set(reflect.Zero(field.Type()))
				decodeInto(r, field.Addr().Interface())

			default:
				return bytesRead, fmt.Errorf("Unhandled Field Kind: %s", field.Kind())
			}
		}
	}

	return bytesRead, nil
}
