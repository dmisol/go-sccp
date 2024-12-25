package sccp

import (
	"fmt"
	"github.com/wmnsk/go-sccp/params"
)

/*
	RLSD (released)
	Table 6/Q.713 − Message type: Released

	Parameter						Clause		Type (F V O)	Length (octets)
	-----------------------------------------------------------------------------
	Message type					2.1				F					1
	Destination local reference		3.2				F					3
	Source local reference			3.3				F					3
	Release cause					3.11			F					1
	Data							3.16			O					3-130
	Importance						3.19			O					3
	End of optional parameter		3.1				O					1
*/

var rlsCauses map[uint8]string

func init() {
	rlsCauses = make(map[uint8]string)
	rlsCauses[0x0] = "End user originated"
	rlsCauses[0x1] = "End user congestion"
	rlsCauses[0x2] = "End user failure"
	rlsCauses[0x3] = "SCCP user originated"
	rlsCauses[0x4] = "Remote procedure error"
	rlsCauses[0x5] = "Inconsistent connection data"
	rlsCauses[0x6] = "Access failure"
	rlsCauses[0x7] = "Access congestion"
	rlsCauses[0x8] = "Subsystem failure"
	rlsCauses[0x9] = "Subsystem congestion"
	rlsCauses[0x0a] = "MTP failure"
	rlsCauses[0x0b] = "Network congestion"
	rlsCauses[0x0c] = "Expiration of reset timer"
	rlsCauses[0x0d] = "Expiration of receive inactivity timer"
	rlsCauses[0x0e] = "Reserved"
	rlsCauses[0x0f] = "Unqualified"
	rlsCauses[0x10] = "SCCP failure"
	// 00010001 - 11110011  are reserved for international use
}

type RLSD struct {
	// Mand:
	Type		MsgType
	DstLocRef	*params.PartyAddress
	SrcLocRef	*params.PartyAddress
	RlsCause	uint8
	// Opt:
	Data		*params.Optional
	Importance	*params.Optional
	EOOptParam	*params.Optional
}

func (m RLSD) RlsCauseName() string {
	name, ok := rlsCauses[m.RlsCause]
	if !ok {
		return "Other"
	}
	return name
}

func ParseRLSD(b []byte) (*RLSD, error) {
	msg := new(RLSD)
	// Type
	msg.Type = MsgType(b[0])
	if msg.Type != MsgTypeRLSD {
		return nil, fmt.Errorf("is not a RLSD message")
	}
	// @todo: Destination local reference
	// ...
	// @todo: Source local reference
	// ...
	// RlsCause
	msg.RlsCause = b[7]
	return msg, nil
}

// @todo
// func (msg RLSD) UnmarshalBinary(b []byte) error {
// }

// @todo
// func (msg RLSD) MarshalBinary() ([]byte, error) {
// }

// @todo
// func (msg RLSD) MarshalLen() int {
// }

// @todo
// func (msg RLSD) MarshalTo(b []byte) error {
// }

// @todo
// func (msg RLSD) String() string {
// }

// @todo
// func (msg RLSD) parseOptional(b []byte) error {
// }

func (msg RLSD) MessageType() MsgType {
	return msg.Type
}

func (msg RLSD) MessageTypeName() string {
	return "RLSD"
}