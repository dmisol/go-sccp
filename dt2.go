package sccp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/dmisol/go-sccp/utils"
)

func NewDT2(dlr uint32, data []byte) *DT2 {
	DT2 := &DT2{
		Type:                      MsgTypeDT2,
		DestinationLocalReference: dlr,
		Sequencing:                0,
		Data:                      data,
	}
	return DT2
}

type DT2 struct {
	Type                      MsgType
	DestinationLocalReference uint32
	Sequencing                uint16 // just dummy for me. All I need is MsgType infact
	Data                      []byte
}

func ParseDT2(b []byte) (*DT2, error) {
	msg := &DT2{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *DT2) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l <= (1 + 3 + 2 + 1) {
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	msg.DestinationLocalReference = utils.Uint24To32(b[1:4])

	msg.Sequencing = binary.BigEndian.Uint16(b[4:])

	if b[6] != 1 { // pointer to var, ae next position
		return io.ErrUnexpectedEOF
	}

	dlen := b[7]
	if l != (dlen + 7 + 1) {
		return io.ErrUnexpectedEOF
	}

	msg.Data = b[8:]
	return nil
}

func (msg *DT2) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *DT2) MarshalLen() int {
	return len(msg.Data) + 8
}

func (msg *DT2) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	copy(b[1:4], utils.Uint32To24(msg.DestinationLocalReference))
	binary.BigEndian.PutUint16(b[4:], msg.Sequencing)
	b[6] = 1
	b[7] = byte(len(msg.Data))
	copy(b[8:], msg.Data)
	return nil
}

func (msg *DT2) String() string {
	return fmt.Sprintf("{Type: DT2, DataLength: %d, Data: %s}", len(msg.Data), hex.EncodeToString(msg.Data))
}

// MessageType returns the Message Type in int.
func (msg *DT2) MessageType() MsgType {
	return msg.Type
}

func (msg *DT2) MessageTypeName() string {
	return "DT2"
}
