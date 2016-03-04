package sflow

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fstelzer/sflow/records"
	"io"
)

type FlowSample struct {
	SequenceNum      uint32
	SourceIdType     byte
	SourceIdIndexVal uint32 // NOTE: this is 3 bytes in the datagram
	SamplingRate     uint32
	SamplePool       uint32
	Drops            uint32
	Input            uint32
	Output           uint32
	numRecords       uint32
	Records          []records.Record
}

func (s FlowSample) String() string {
	type X FlowSample
	x := X(s)
	return fmt.Sprintf("FlowSample: %+v", x)
}

// SampleType returns the type of sFlow sample.
func (s *FlowSample) SampleType() int {
	return TypeFlowSample
}

func (s *FlowSample) GetRecords() []records.Record {
	return s.Records
}

func decodeFlowSample(r io.ReadSeeker) (Sample, error) {
	s := &FlowSample{}

	var err error

	err = binary.Read(r, binary.BigEndian, &s.SequenceNum)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.SourceIdType)
	if err != nil {
		return nil, err
	}

	var srcIdIndexVal [3]byte
	n, err := r.Read(srcIdIndexVal[:])
	if err != nil {
		return nil, err
	}

	if n != 3 {
		return nil, errors.New("sflow: counter sample decoding error")
	}

	s.SourceIdIndexVal = uint32(srcIdIndexVal[2]) |
		uint32(srcIdIndexVal[1]<<8) |
		uint32(srcIdIndexVal[0]<<16)

	err = binary.Read(r, binary.BigEndian, &s.SamplingRate)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.SamplePool)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.Drops)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.Input)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.Output)
	if err != nil {
		return nil, err
	}

	err = binary.Read(r, binary.BigEndian, &s.numRecords)
	if err != nil {
		return nil, err
	}

	for i := uint32(0); i < s.numRecords; i++ {
		format, length := uint32(0), uint32(0)

		err = binary.Read(r, binary.BigEndian, &format)
		if err != nil {
			return nil, err
		}

		err = binary.Read(r, binary.BigEndian, &length)
		if err != nil {
			return nil, err
		}

		var rec records.Record

		if rec, err = records.DecodeFlow(r, format); err != nil {
			fmt.Printf("Error: %s\n", err)
			_, err := r.Seek(int64(length), 1)
			if err != nil {
				return nil, err
			}
			continue
		}

		s.Records = append(s.Records, rec)
	}

	return s, nil
}

func (s *FlowSample) encode(w io.Writer) error {
	var err error

	// We first need to encode the records.
	buf := &bytes.Buffer{}

	for _, rec := range s.Records {
		err = rec.Encode(buf)
		if err != nil {
			return records.ErrEncodingRecord
		}
	}

	// Fields
	encodedSampleSize := uint32(4 + 1 + 3 + 4 + 4 + 4 + 4 + 4 + 4)

	// Encoded records
	encodedSampleSize += uint32(buf.Len())

	err = binary.Write(w, binary.BigEndian, uint32(s.SampleType()))
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, encodedSampleSize)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.SequenceNum)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian,
		uint32(s.SourceIdType)|s.SourceIdIndexVal<<24)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.SamplingRate)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.SamplePool)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.Drops)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.Input)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, s.Output)
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.BigEndian, uint32(len(s.Records)))
	if err != nil {
		return err
	}

	_, err = io.Copy(w, buf)
	return err
}
