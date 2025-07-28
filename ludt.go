package sccp

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/dmisol/go-sccp/params"
)

type LUDT struct {
	Type MsgType
	params.ProtocolClass
	HopCntr             uint8
	ptrs                [4]uint16
	CalledPartyAddress  *params.PartyAddress
	CallingPartyAddress *params.PartyAddress

	DataLength uint16

	Data []byte

	// fixme: will not use it/will not rx this. Fell free to implement urself
	// Segmentation
	// Importance
}

func NewLUDT(pcls int, mhandle bool, cdpa, cgpa *params.PartyAddress, data []byte) *LUDT {
	u := &LUDT{
		Type: MsgTypeLUDT,
		ProtocolClass: params.NewProtocolClass(
			pcls, mhandle,
		),
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
		Data:                data,
	}
	u.DataLength = uint16(len(data))
	u.ptrs[0] = 8
	u.ptrs[1] = u.ptrs[0] + uint16(cdpa.Length) - 2
	u.ptrs[2] = u.ptrs[1] + uint16(cgpa.Length) - 2
	u.ptrs[3] = u.ptrs[2] + u.DataLength

	return u
}

// MarshalBinary returns the byte sequence generated from a LUDT instance.
func (u *LUDT) MarshalBinary() ([]byte, error) {
	fmt.Println("c")

	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (u *LUDT) MarshalTo(b []byte) error {
	l := uint16(len(b))
	if l < 12 {
		return io.ErrUnexpectedEOF
	}

	if u.ptrs[3]+5 >= l {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(u.Type)
	b[1] = uint8(u.ProtocolClass)
	b[2] = u.HopCntr

	for i := 0; i < 4; i++ {
		binary.BigEndian.PutUint16(b[3+2*i:], u.ptrs[i])
	}

	fmt.Println("b")
	if err := u.CalledPartyAddress.MarshalTo(b[7:int(u.ptrs[1]+4)]); err != nil {
		return err
	}
	if err := u.CallingPartyAddress.MarshalTo(b[int(u.ptrs[1]+4):int(u.ptrs[2]+5)]); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(b[int(u.ptrs[2]+5):], u.DataLength)
	copy(b[int(u.ptrs[2]+6):int(u.ptrs[3]+7)], u.Data)
	b[int(u.ptrs[3]+7)] = 0

	return nil
}

// ParseLUDT decodes given byte sequence as a SCCP LUDT.
func ParseLUDT(b []byte) (*LUDT, error) {
	u := &LUDT{}
	if err := u.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP LUDT.
func (u *LUDT) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l <= 14 {
		return io.ErrUnexpectedEOF
	}

	u.Type = MsgType(b[0])
	u.ProtocolClass = params.ProtocolClass(b[1])
	u.HopCntr = b[2]
	for i := 0; i < 4; i++ {
		u.ptrs[i] = binary.BigEndian.Uint16(b[3+2*i:])
	}
	fmt.Println(u.ptrs)
	if l < 3+int(u.ptrs[0])+1 {
		return io.ErrUnexpectedEOF
	}
	if l < 4+int(u.ptrs[1])+1 {
		return io.ErrUnexpectedEOF
	}
	if l < 5+int(u.ptrs[2])+1 {
		return io.ErrUnexpectedEOF
	}
	if l < 6+int(u.ptrs[3]) { // where u.Data starts
		return io.ErrUnexpectedEOF
	}

	var err error
	u.CalledPartyAddress, err = params.ParsePartyAddress(b[11 : int(u.ptrs[1])+5])
	if err != nil {
		return err
	}

	u.CallingPartyAddress, err = params.ParsePartyAddress(b[int(u.ptrs[1]+5):int(u.ptrs[2]+7)])
	if err != nil {
		return err
	}

	u.DataLength = binary.BigEndian.Uint16(b[int(u.ptrs[2]+7):])

	if offset, dataLen := int(u.ptrs[2]+9), int(u.DataLength); l >= offset+dataLen {
		u.Data = b[offset : offset+dataLen]
		return nil
	}

	return io.ErrUnexpectedEOF
}

// MarshalLen returns the serial length.
func (u *LUDT) MarshalLen() int {
	l := 6 + 3
	if param := u.CalledPartyAddress; param != nil {
		l += param.MarshalLen()
	}
	if param := u.CallingPartyAddress; param != nil {
		l += param.MarshalLen()
	}
	l += len(u.Data)

	return l
}

// String returns the LUDT values in human readable format.
func (u *LUDT) String() string {
	return fmt.Sprintf("{Type: %d, CalledPartyAddress: %v, CallingPartyAddress: %v, DataLength: %d, Data: %x}",
		u.Type,
		u.CalledPartyAddress,
		u.CallingPartyAddress,
		u.DataLength,
		u.Data,
	)
}

// MessageType returns the Message Type in int.
func (u *LUDT) MessageType() MsgType {
	return MsgTypeLUDT
}

// MessageTypeName returns the Message Type in string.
func (u *LUDT) MessageTypeName() string {
	return "LUDT"
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (u *LUDT) CdGT() string {
	return u.CalledPartyAddress.GTString()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (u *LUDT) CgGT() string {
	return u.CallingPartyAddress.GTString()
}
