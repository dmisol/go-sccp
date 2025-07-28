package sccp

import (
	"bytes"
	"fmt"
	"testing"
)

var testXUDT = []byte{byte(MsgTypeXUDT), 0x0,
	0x0,     // hop cntr, inserted
	0x3 + 1, //called, increased (poiner added)
	0x5 + 1, // calling,...
	0x7 + 1, // data
	29,      // ptr to optional, inserted
	0x2, 0x42, 0xfe,
	0x2, 0x42, 0xfe,
	0x16,
	0x0, 0x14, 0x52, 0x8, 0x8, 0x29, 0x5, 0x70, 0x48, 0x96, 0x10, 0x14, 0x72, 0x9, 0x4, 0x20, 0x2, 0xf8, 0xc5, 0x1a, 0x1, 0x6,
	0x0, // added end of optional
}

func TestXUDT(t *testing.T) {
	t.Skip()

	x, err := ParseXUDT(testXUDT)
	if err != nil {
		t.Fatal(err)
	}

	y, err := x.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	cmp := bytes.Compare(testXUDT, y)
	if cmp != 0 {
		fmt.Println(testXUDT)
		fmt.Println(y)
		t.Fail()
	}
}
