package types

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *StorerDocument) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxvk uint32
	zxvk, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxvk > 0 {
		zxvk--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.DocId, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "lat":
			z.Latitude, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "long":
			z.Longitude, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		case "f":
			z.Fields, err = dc.ReadIntf()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *StorerDocument) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "id"
	err = en.Append(0x84, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteUint64(z.DocId)
	if err != nil {
		return
	}
	// write "lat"
	err = en.Append(0xa3, 0x6c, 0x61, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Latitude)
	if err != nil {
		return
	}
	// write "long"
	err = en.Append(0xa4, 0x6c, 0x6f, 0x6e, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteFloat64(z.Longitude)
	if err != nil {
		return
	}
	// write "f"
	err = en.Append(0xa1, 0x66)
	if err != nil {
		return err
	}
	err = en.WriteIntf(z.Fields)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *StorerDocument) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "id"
	o = append(o, 0x84, 0xa2, 0x69, 0x64)
	o = msgp.AppendUint64(o, z.DocId)
	// string "lat"
	o = append(o, 0xa3, 0x6c, 0x61, 0x74)
	o = msgp.AppendFloat64(o, z.Latitude)
	// string "long"
	o = append(o, 0xa4, 0x6c, 0x6f, 0x6e, 0x67)
	o = msgp.AppendFloat64(o, z.Longitude)
	// string "f"
	o = append(o, 0xa1, 0x66)
	o, err = msgp.AppendIntf(o, z.Fields)
	if err != nil {
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *StorerDocument) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.DocId, bts, err = msgp.ReadUint64Bytes(bts)
			if err != nil {
				return
			}
		case "lat":
			z.Latitude, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "long":
			z.Longitude, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		case "f":
			z.Fields, bts, err = msgp.ReadIntfBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *StorerDocument) Msgsize() (s int) {
	s = 1 + 3 + msgp.Uint64Size + 4 + msgp.Float64Size + 5 + msgp.Float64Size + 2 + msgp.GuessSize(z.Fields)
	return
}
