package message

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *MessageParseError) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var tmp string
		tmp, err = dc.ReadString()
		(*z) = MessageParseError(tmp)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z MessageParseError) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteString(string(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MessageParseError) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendString(o, string(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageParseError) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var tmp string
		tmp, bts, err = msgp.ReadStringBytes(bts)
		(*z) = MessageParseError(tmp)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

func (z MessageParseError) Msgsize() (s int) {
	s = msgp.StringPrefixSize + len(string(z))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *UDPMessage) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "context":
			z.Context, err = dc.ReadString()
			if err != nil {
				return
			}
		case "name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "type":
			{
				var tmp int
				tmp, err = dc.ReadInt()
				z.Type = MessageType(tmp)
			}
			if err != nil {
				return
			}
		case "time":
			z.Time, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "host":
			z.Host, err = dc.ReadString()
			if err != nil {
				return
			}
		case "autoIP":
			z.FillHostWithIP, err = dc.ReadBool()
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
func (z *UDPMessage) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	err = en.Append(0x86)
	if err != nil {
		return err
	}
	// write "context"
	err = en.Append(0xa7, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Context)
	if err != nil {
		return
	}
	// write "name"
	err = en.Append(0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "type"
	err = en.Append(0xa4, 0x74, 0x79, 0x70, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt(int(z.Type))
	if err != nil {
		return
	}
	// write "time"
	err = en.Append(0xa4, 0x74, 0x69, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteTime(z.Time)
	if err != nil {
		return
	}
	// write "host"
	err = en.Append(0xa4, 0x68, 0x6f, 0x73, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Host)
	if err != nil {
		return
	}
	// write "autoIP"
	err = en.Append(0xa6, 0x61, 0x75, 0x74, 0x6f, 0x49, 0x50)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.FillHostWithIP)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *UDPMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	o = append(o, 0x86)
	// string "context"
	o = append(o, 0xa7, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74)
	o = msgp.AppendString(o, z.Context)
	// string "name"
	o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.Name)
	// string "type"
	o = append(o, 0xa4, 0x74, 0x79, 0x70, 0x65)
	o = msgp.AppendInt(o, int(z.Type))
	// string "time"
	o = append(o, 0xa4, 0x74, 0x69, 0x6d, 0x65)
	o = msgp.AppendTime(o, z.Time)
	// string "host"
	o = append(o, 0xa4, 0x68, 0x6f, 0x73, 0x74)
	o = msgp.AppendString(o, z.Host)
	// string "autoIP"
	o = append(o, 0xa6, 0x61, 0x75, 0x74, 0x6f, 0x49, 0x50)
	o = msgp.AppendBool(o, z.FillHostWithIP)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *UDPMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var isz uint32
	isz, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for isz > 0 {
		isz--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "context":
			z.Context, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "name":
			z.Name, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "type":
			{
				var tmp int
				tmp, bts, err = msgp.ReadIntBytes(bts)
				z.Type = MessageType(tmp)
			}
			if err != nil {
				return
			}
		case "time":
			z.Time, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "host":
			z.Host, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "autoIP":
			z.FillHostWithIP, bts, err = msgp.ReadBoolBytes(bts)
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

func (z *UDPMessage) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.Context) + 5 + msgp.StringPrefixSize + len(z.Name) + 5 + msgp.IntSize + 5 + msgp.TimeSize + 5 + msgp.StringPrefixSize + len(z.Host) + 7 + msgp.BoolSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MessageType) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var tmp int
		tmp, err = dc.ReadInt()
		(*z) = MessageType(tmp)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z MessageType) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteInt(int(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z MessageType) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendInt(o, int(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MessageType) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var tmp int
		tmp, bts, err = msgp.ReadIntBytes(bts)
		(*z) = MessageType(tmp)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

func (z MessageType) Msgsize() (s int) {
	s = msgp.IntSize
	return
}
