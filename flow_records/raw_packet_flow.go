package flow_records

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
)

// RawPacketFlow is a raw Ethernet header flow record.
type RawPacketFlow struct {
	Protocol      uint32
	FrameLength   uint32
	Stripped      uint32
	HeaderSize    uint32
	Header        []byte
	DecodedHeader map[string]interface{}
}

type EthernetHeader struct {
	DstMac net.HardwareAddr
	SrcMac net.HardwareAddr
}

/* define my own IP header struct - to ease portability */
type IPv4Header struct {
	VersionAndLen uint8
	Tos           uint8
	TotLen        uint16
	Id            uint16
	FragOff       uint16
	Ttl           uint8
	Protocol      uint8
	Check         uint16
	SrcAddr       net.IP `ipVersion:"4"`
	DstAddr       net.IP `ipVersion:"4"`
}

type IPv6Header struct {
	VersionAndPriority uint8
	Label1             uint8
	Label2             uint8
	Label3             uint8
	PayloadLength      uint16
	NextHeader         uint8
	Ttl                uint8
	SrcAddr            net.IP `ipVersion:"6"`
	DstAddr            net.IP `ipVersion:"6"`
}

type TcpHeader struct {
	SrcPort  uint16
	DstPort  uint16
	Seq      uint32
	Ack      uint32
	UnUsed   uint8
	Flags    uint8
	Window   uint16
	Checksum uint16
	Urgent   uint16
}

type UdpHeader struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
}

type IcmpHeader struct {
	Type uint8
	Code uint8
}

func (f RawPacketFlow) String() string {
	type X RawPacketFlow
	x := X(f)
	return fmt.Sprintf("RawPacketFlow: %+v", x)
}

// RecordType returns the type of flow record.
func (f RawPacketFlow) RecordType() int {
	return TypeRawPacketFlowRecord
}

func (f *RawPacketFlow) decodeIPHeader(ipVersion int, h io.Reader) error {
	var err error

	if ipVersion == 4 {
		ip := IPv4Header{}

		if err = Decode(h, &ip); err != nil {
			fmt.Printf("Err: %s\n", err)
			return err
		}
		f.DecodedHeader["ip"] = ip

		//Can we decode a following Layer4 Protocol Header?
		switch ip.Protocol {
		case IPProtocolTCP:
			tcp := TcpHeader{}
			if err = Decode(h, &tcp); err != nil {
				return err
			}
			f.DecodedHeader["tcp"] = tcp
		case IPProtocolUDP:
			udp := UdpHeader{}
			if err = Decode(h, &udp); err != nil {
				return err
			}
			f.DecodedHeader["udp"] = udp
		case IPProtocolICMP:
			icmp := IcmpHeader{}
			if err = Decode(h, &icmp); err != nil {
				return err
			}
			f.DecodedHeader["icmp"] = icmp
		default:
			fmt.Printf("Unknown Protocol: %d\n", ip.Protocol)
		}

	} else if ipVersion == 6 {
		//FIXME: IPv6 has complex Header Extensions
		//FIXME: IMPLEMENT ME
		return fmt.Errorf("IPv6 is not implemented yet")

		IPHeader := IPv6Header{}

		if err = Decode(h, &IPHeader); err != nil {
			return err
		}
		f.DecodedHeader["ip"] = IPHeader
	}

	return nil
}

func (f *RawPacketFlow) decodeHeader(headerType uint32) error {
	var err error = nil

	f.DecodedHeader = make(map[string]interface{})

	if len(f.Header) < MinimumEthernetHeaderSize {
		return nil
	}

	h := bytes.NewReader(f.Header)

	switch headerType {
	case HeaderProtocolEthernetISO8023:
		ethernet := EthernetHeader{}
		if err = Decode(h, &ethernet); err != nil {
			return err
		}
		f.DecodedHeader["ethernet"] = ethernet

		// Determine the Type of the next Header
		buffer := make([]byte, 2)
		if err = binary.Read(h, binary.BigEndian, &buffer); err != nil {
			fmt.Printf("Error: %s\n", err)
			return err
		}

		//TODO: Handle VSNAP / 802.2/802 &  IPX

		switch hex.EncodeToString(buffer) {
		case HeaderTypeIPv4:
			if err = f.decodeIPHeader(4, h); err != nil {
				return err
			}
		case HeaderTypeIPv6:
			if err = f.decodeIPHeader(6, h); err != nil {
				return err
			}
		}
	case HeaderProtocolIPv4:
		if err = f.decodeIPHeader(4, h); err != nil {
			return err
		}
	case HeaderProtocolIPv6:
		if err = f.decodeIPHeader(6, h); err != nil {
			return err
		}
	default:
		fmt.Printf("Unknown Headertype: %d\n", headerType)
	}

	//fmt.Printf("Headers: %+#v\n", f.DecodedHeader)
	return err
}

func DecodeRawPacketFlow(r io.Reader) (RawPacketFlow, error) {
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

	// Try to decode the retrieved headers
	if err = f.decodeHeader(f.Protocol); err != nil {
		return f, err
	}

	return f, err
}

func (f RawPacketFlow) Encode(w io.Writer) error {
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

	// We don't need to reencode the DecodedHeaders as the raw data is still in the Header Field.

	return err
}
