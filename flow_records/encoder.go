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