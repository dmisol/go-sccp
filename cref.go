package sccp

import (
	"fmt"
	"io"

	"github.com/dmisol/go-sccp/params"
	"github.com/dmisol/go-sccp/utils"
)

func NewCREF(dlr uint32, cause byte, opts []*params.Optional) *CREF {
	CREF := &CREF{
		Type:                      MsgTypeCREF,
		DestinationLocalReference: dlr,
		Cause:                     cause,
		Opts:                      opts,
	}
	return CREF
}

type CREF struct {
	Type                      MsgType
	DestinationLocalReference uint32
	Cause                     byte

	Opts []*params.Optional

	Data               *params.Optional
	CalledPartyAddress *params.PartyAddress
}

func ParseCREF(b []byte) (*CREF, error) {
	msg := &CREF{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *CREF) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l < 6 {
		return io.ErrUnexpectedEOF
	}
	msg.Type = MsgType(b[0])
	msg.DestinationLocalReference = utils.Uint24To32(b[1:4])
	msg.Cause = b[4]

	optr := b[5]
	if optr == 0 {
		return nil
	}
	if optr != 1 {
		return io.ErrUnexpectedEOF
	}

	if err := msg.parseOptional(b[6:]); err != nil {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (msg *CREF) parseOptional(b []byte) error {
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

		switch t {
		case params.DataTag:
			msg.Data = o
		case params.CdPtyAddrTag:
			var err error
			msg.CalledPartyAddress, err = params.ParsePartyAddress(b[p : p+2+l])
			if err != nil {
				return err
			}
		}

		msg.Opts = append(msg.Opts, o)
		p += 2 + l
	}
	return nil
}

// MarshalBinary returns the byte sequence generated from a CREF instance.
func (msg *CREF) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *CREF) MarshalLen() int {
	if len(msg.Opts) == 0 {
		return 6 // 5 fixed + 1 ptr
	}
	l := 7 // 5 fixed + 1 ptr + last optional
	for _, v := range msg.Opts {
		l += int(v.Len) + 2
	}

	return l
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (msg *CREF) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	copy(b[1:4], utils.Uint32To24(msg.DestinationLocalReference))
	b[4] = byte(msg.Cause)

	if len(msg.Opts) == 0 {
		return nil
	}

	b[5] = 1
	p := uint8(6)

	for i := 0; i < len(msg.Opts); i++ {
		b[p] = msg.Opts[i].Tag
		b[p+1] = msg.Opts[i].Len
		copy(b[p+2:], msg.Opts[i].Value)

		p += msg.Opts[i].Len + 2
	}
	return nil
}

func (msg *CREF) String() string {
	if msg.CalledPartyAddress != nil {
		return fmt.Sprintf("{Type: CREF, CalledPartyAddress: %v}", msg.CalledPartyAddress)
	}
	return "{Type: CREF}"
}

// MessageType returns the Message Type in int.
func (msg *CREF) MessageType() MsgType {
	return msg.Type
}

func (msg *CREF) MessageTypeName() string {
	return "CREF"
}
