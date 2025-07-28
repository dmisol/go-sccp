package sccp

import (
	"fmt"
	"io"

	"github.com/dmisol/go-sccp/params"
)

type XUDT struct {
	Type MsgType
	params.ProtocolClass
	HopCntr             uint8
	ptrs                [4]uint8
	CalledPartyAddress  *params.PartyAddress
	CallingPartyAddress *params.PartyAddress

	DataLength uint8 // have to inherit it from the original git :-(

	Data []byte

	// fixme: will not use it/will not rx this. Fell free to implement urself
	// Segmentation
	// Importance
}

func NewXUDT(pcls int, mhandle bool, cdpa, cgpa *params.PartyAddress, data []byte) *XUDT {
	u := &XUDT{
		Type: MsgTypeXUDT,
		ProtocolClass: params.NewProtocolClass(
			pcls, mhandle,
		),
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
		Data:                data,
	}
	u.DataLength = uint8(len(data))
	u.ptrs[0] = 4
	u.ptrs[1] = u.ptrs[0] + cdpa.Length
	u.ptrs[2] = u.ptrs[1] + cgpa.Length
	u.ptrs[3] = u.ptrs[2] + u.DataLength

	return u
}

// MarshalBinary returns the byte sequence generated from a XUDT instance.
func (u *XUDT) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (u *XUDT) MarshalTo(b []byte) error {
	l := uint8(len(b))
	if l < 8 {
		return io.ErrUnexpectedEOF
	}

	if u.ptrs[3]+5 >= l {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(u.Type)
	b[1] = uint8(u.ProtocolClass)
	b[2] = u.HopCntr
	copy(b[3:], u.ptrs[0:])

	if err := u.CalledPartyAddress.MarshalTo(b[7:int(u.ptrs[1]+4)]); err != nil {
		return err
	}
	if err := u.CallingPartyAddress.MarshalTo(b[int(u.ptrs[1]+4):int(u.ptrs[2]+5)]); err != nil {
		return err
	}
	b[int(u.ptrs[2]+5)] = u.DataLength
	copy(b[int(u.ptrs[2]+6):int(u.ptrs[3]+7)], u.Data)
	b[int(u.ptrs[3]+7)] = 0

	return nil
}

// ParseXUDT decodes given byte sequence as a SCCP XUDT.
func ParseXUDT(b []byte) (*XUDT, error) {
	u := &XUDT{}
	if err := u.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP XUDT.
func (u *XUDT) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l <= 6 { // where CdPA starts
		return io.ErrUnexpectedEOF
	}

	u.Type = MsgType(b[0])
	u.ProtocolClass = params.ProtocolClass(b[1])
	u.HopCntr = b[2]
	copy(u.ptrs[0:], b[3:])
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
	u.CalledPartyAddress, err = params.ParsePartyAddress(b[7:int(u.ptrs[1]+4)])
	if err != nil {
		return err
	}

	u.CallingPartyAddress, err = params.ParsePartyAddress(b[int(u.ptrs[1]+4):int(u.ptrs[2]+5)])
	if err != nil {
		return err
	}

	// succeed if the rest of buffer is longer than u.DataLength
	u.DataLength = b[int(u.ptrs[2]+5)]
	if offset, dataLen := int(u.ptrs[2]+6), int(u.DataLength); l >= offset+dataLen {
		u.Data = b[offset : offset+dataLen]
		return nil
	}

	return io.ErrUnexpectedEOF
}

// MarshalLen returns the serial length.
func (u *XUDT) MarshalLen() int {
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

// String returns the XUDT values in human readable format.
func (u *XUDT) String() string {
	return fmt.Sprintf("{Type: %d, CalledPartyAddress: %v, CallingPartyAddress: %v, DataLength: %d, Data: %x}",
		u.Type,
		u.CalledPartyAddress,
		u.CallingPartyAddress,
		u.DataLength,
		u.Data,
	)
}

// MessageType returns the Message Type in int.
func (u *XUDT) MessageType() MsgType {
	return MsgTypeXUDT
}

// MessageTypeName returns the Message Type in string.
func (u *XUDT) MessageTypeName() string {
	return "XUDT"
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (u *XUDT) CdGT() string {
	return u.CalledPartyAddress.GTString()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (u *XUDT) CgGT() string {
	return u.CallingPartyAddress.GTString()
}
