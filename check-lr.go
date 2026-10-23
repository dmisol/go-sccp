package sccp

import "github.com/dmisol/go-sccp/utils"

// used for SCCP-based load balancer
func DLR(b []byte) (uint32, bool) {
	if len(b) < 4 {
		return 0, false
	}

	switch MsgType(b[0]) {
	case MsgTypeCC,
		MsgTypeDT1, MsgTypeDT2,
		MsgTypeRLC, MsgTypeRLSD,
		MsgTypeCREF, MsgTypeRSR, MsgTypeRSC,
		MsgTypeERR, MsgTypeIT:
		return utils.Uint24To32(b[1:4]), true
	}

	return 0, false
}

func SLR(b []byte) (uint32, bool) {
	if len(b) < 4 {
		return 0, false
	}

	switch MsgType(b[0]) {
	case MsgTypeCC,
		MsgTypeDT1, MsgTypeDT2,
		MsgTypeRLC, MsgTypeRLSD,
		MsgTypeCREF, MsgTypeRSR, MsgTypeRSC,
		MsgTypeERR, MsgTypeIT:
		return utils.Uint24To32(b[4:7]), true
	}

	return 0, false
}
