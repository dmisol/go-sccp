package sccp

import (
	"fmt"
	"io"
	"github.com/wmnsk/go-sccp/params"
)

/*
	RLSD (released)
	Table 6/Q.713 − Message type: Released

	The RLSD message contains:
	- one pointer;
	- the parameters indicated in Table 6.

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
	// 00010001 - 11110011  are reserved for the international use
}

type RLSD struct {
	// Mand:
	Type		MsgType
	DstLocRef	params.LocalReference
	SrcLocRef	params.LocalReference
	RlsCause	uint8
	// Opt:
	Data		*params.Optional
	Importance	*params.Optional // 0-7
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
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return msg, nil
}

// @todo
func (msg *RLSD) UnmarshalBinary(b []byte) error {
	l := len(b)
	var minLenB int = 0 +
		1 + // Message type
		3 + // Destination local reference
		3 + // Source local reference
		1 +	// Release cause
		1	// Pointer to the optional parameter
	if l <= minLenB {
		return io.ErrUnexpectedEOF
	}

	// MANDATORY:

	// Mand: Message type
	msg.Type = MsgType(b[0])
	if msg.Type != MsgTypeRLSD {
		return fmt.Errorf("is not a RLSD message")
	}

	// Mand: Destination local reference
	if err := msg.DstLocRef.Read(b[1:4]); err != nil {
		return err
	}

	// Mand: Source local reference
	if err := msg.SrcLocRef.Read(b[4:7]); err != nil {
		return err
	}

	// Mand: Release cause
	msg.RlsCause = b[7]

	// OPTIONAL:

	var byteIdx int = 8
	optPrmsPtr := int(b[byteIdx])

	if optPrmsPtr == 0x0 { // No optional parameters
		return nil
	}

	opts, err := ParseOptionalParameters(b[byteIdx+optPrmsPtr:])
	if err != nil {
		return err
	}

	// Opt: data
	if opt, ok := opts[params.DataTag]; ok {
		msg.Data = opt
	}

	// Opt: importance
	if opt, ok := opts[params.ImportanceTag]; ok {
		msg.Importance = opt
	}

	return nil
}

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

func (msg RLSD) MessageType() MsgType {
	return msg.Type
}

func (msg RLSD) MessageTypeName() string {
	return "RLSD"
}