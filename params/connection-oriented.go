package params

import "fmt"

const (
	CgPtyAddrTag			uint8 = 0x04
	ProtocolClassTag		uint8 = 0x05
	// @todo: Segmenting/reassembling
	// @todo: Receive sequence number
	// @todo: Sequencing/segmenting
	CreditTag				uint8 = 0x09
	ReleaseCauseTag			uint8 = 0x0A
	// @todo: Return cause
	// @todo: Reset cause
	// @todo: Error cause
	// @todo: Refusal cause
	DataTag					uint8 = 0x0F
	// @todo: Segmentation
	HopCounterTag 			uint8 = 0x11
	ImportanceTag 			uint8 = 0x12
	// @todo: Long data
)

type LocalReference struct {
	Value uint32
} // just 24 bits used
func (lr *LocalReference) Read(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("unable to read local reference: given bytes length is invalid")
	}
	b[0] = byte((lr.Value >> 16) & 0xFF)
	b[1] = byte((lr.Value >> 8) & 0xFF)
	b[2] = byte(lr.Value & 0xFF)
	return nil
}
func (lr *LocalReference) Write(b [3]byte) {
	lr.Value = (uint32(b[0])<<16)&0xFF0000 + (uint32(b[1])<<8)&0xFF00 + uint32(b[2])&0xFF
}
func (lr *LocalReference) String() string {
	return fmt.Sprintf("%05X", lr.Value)
}

type Optional struct {
	Tag   uint8
	Len   uint8
	Value []byte
	Ptr   uint8
}
