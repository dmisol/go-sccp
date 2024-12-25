package sccp

import "github.com/wmnsk/go-sccp/params"

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
	rlsCauses[0b00000000] = "End user originated"
	rlsCauses[0b00000001] = "End user congestion"
	rlsCauses[0b00000010] = "End user failure"
	rlsCauses[0b00000011] = "SCCP user originated"
	rlsCauses[0b00000100] = "Remote procedure error"
	rlsCauses[0b00000101] = "Inconsistent connection data"
	rlsCauses[0b00000110] = "Access failure"
	rlsCauses[0b00000111] = "Access congestion"
	rlsCauses[0b00001000] = "Subsystem failure"
	rlsCauses[0b00001001] = "Subsystem congestion"
	rlsCauses[0b00001010] = "MTP failure"
	rlsCauses[0b00001011] = "Network congestion"
	rlsCauses[0b00001100] = "Expiration of reset timer"
	rlsCauses[0b00001101] = "Expiration of receive inactivity timer"
	rlsCauses[0b00001110] = "Reserved"
	rlsCauses[0b00001111] = "Unqualified"
	rlsCauses[0b00010000] = "SCCP failure"
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