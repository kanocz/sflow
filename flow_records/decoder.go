package flow_records

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

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
							Decode(r, field.Index(x).Addr().Interface())
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
