package sccp

import (
	"io"

	"github.com/dmisol/go-sccp/params"
	"github.com/dmisol/go-sccp/utils"
)

func NewRLSD(dlr uint32, slr uint32, cause byte, opts []*params.Optional) *RLSD {
	RLSD := &RLSD{
		Type:                      MsgTypeRLSD,
		DestinationLocalReference: dlr,
		SourceLocalReference:      slr,
	}
	return RLSD
}

type RLSD struct {
	Type                      MsgType
	DestinationLocalReference uint32
	SourceLocalReference      uint32
	Cause                     byte

	Opts []*params.Optional
}

func ParseRLSD(b []byte) (*RLSD, error) {
	msg := &RLSD{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *RLSD) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l < 9 {
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	msg.DestinationLocalReference = utils.Uint24To32(b[1:4])
	msg.SourceLocalReference = utils.Uint24To32(b[4:7])
	msg.Cause = b[7]

	if b[8] == 0 {
		return nil
	}
	if b[8] != 1 {
		return io.ErrUnexpectedEOF
	}
	return msg.parseOptional(b[9:])
}

func (msg *RLSD) parseOptional(b []byte) error {
	p := uint8(0)
	for p < uint8(len(b)) {
		t := b[p]

		if t == 0 {
			return nil
		}
		if (p + 1) >= uint8(len(b)) {
			return io.ErrUnexpectedEOF
		}

		l := b[p+1]
		if (p + 1 + l) >= uint8(len(b)) {
			return io.ErrUnexpectedEOF
		}

		o := &params.Optional{
			Tag:   t,
			Len:   l,
			Value: b[p+2 : p+2+l],
		}

		msg.Opts = append(msg.Opts, o)
		p += 2 + l
	}

	return nil
}

func (msg *RLSD) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *RLSD) MarshalLen() int {
	if len(msg.Opts) == 0 {
		return 9 // 8 fixed + 0 ptr
	}
	l := 10 // 8 fixed + 1 ptr + last optional
	for _, v := range msg.Opts {
		l += int(v.Len) + 2
	}

	return l
}

func (msg *RLSD) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	copy((b[1:4]), utils.Uint32To24(msg.DestinationLocalReference))
	copy(b[4:], utils.Uint32To24(msg.SourceLocalReference))
	b[7] = msg.Cause

	if len(msg.Opts) == 0 {
		return nil
	}

	b[8] = 1
	p := uint8(9)

	for i := 0; i < len(msg.Opts); i++ {
		b[p] = msg.Opts[i].Tag
		b[p+1] = msg.Opts[i].Len
		copy(b[p+2:], msg.Opts[i].Value)

		p += msg.Opts[i].Len + 2
	}
	return nil
}

func (msg *RLSD) String() string {
	return "{Type: RLSD}"
}

// MessageType returns the Message Type in int.
func (msg *RLSD) MessageType() MsgType {
	return msg.Type
}

func (msg *RLSD) MessageTypeName() string {
	return "RLSD"
}
