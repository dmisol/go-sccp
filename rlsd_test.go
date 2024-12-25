package sccp

import (
	"testing"
)

var validRLSDMock []byte = []byte{0x4, 0x0, 0x3, 0x8c, 0x0, 0x60, 0x4e, 0x0d, 0x0, 0x4, 0x0, 0x3, 0x8c, 0x0, 0x60, 0x4e, 0x0, 0x0}
// var invalidRLSDMock []byte = []byte{0x1}

func TestRLSDType(t *testing.T) {
	msg, err := ParseRLSD(validRLSDMock)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Type != MsgTypeRLSD {
		t.Fatal("Msg is not a RLSD")
	}
}

func TestRLSDRlsCause(t *testing.T) {
	msg, _ := ParseRLSD(validRLSDMock)
	if msg.RlsCause != 0xd {
		t.Fatal("Invalid rls cause")
	}
}

func TestRlsCause(t *testing.T) {
	msg, _ := ParseRLSD(validRLSDMock)
	if msg.RlsCauseName() != "Expiration of receive inactivity timer" {
		t.Fatal("Rls cause name was not received")
	}
}

// func TestRLSDRlsCauseName(t *testing.T) {
// 	msg, err := ParseRLSD(validRLSDMock)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if msg.RlsCauseName() != "SCCP user originated" {
// 		t.Fatal("Invalid rls cause name")
// 	}
// }
