package sccp

import (
	"bytes"
	"fmt"
	"testing"
)

var testLUDT = []byte{byte(MsgTypeLUDT), 0x0,
	0x0,  // hop cntr, inserted
	0, 8, // 0x3 + 1, //called, increased (poiner added)
	0, 9, // 0x5 + 1, // calling,...
	0, 10, //0x7 + 1, // data
	0, 32, // ptr to optional, inserted
	0x2, 0x42, 0xfe,
	0x2, 0x42, 0xfe,
	0, 0x16, // int8->int16
	0x0, 0x14, 0x52, 0x8, 0x8, 0x29, 0x5, 0x70, 0x48, 0x96, 0x10, 0x14, 0x72, 0x9, 0x4, 0x20, 0x2, 0xf8, 0xc5, 0x1a, 0x1, 0x6,
	0x0, // added end of optional
}

func TestLUDT(t *testing.T) {

	x, err := ParseLUDT(testLUDT)
	if err != nil {
		t.Fatal(err)
	}

	y, err := x.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	cmp := bytes.Compare(testLUDT, y)
	if cmp != 0 {
		fmt.Println(testLUDT)
		fmt.Println(y)
		t.Fail()
	}
}
