package sccp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

var mockCREFs = [][]byte{
	{0x3, 0x0, 0x3, 0x75, 0x1, 0x0},
	{0x3, 0x0, 0x5f, 0x5, 0x1, 0x1, 0x3, 0x4, 0x43, 0x1c, 0x2d, 0xfe, 0x0},
}

func TestCREF(t *testing.T) {
	for i, v := range mockCREFs {
		cref, err := ParseCREF(v)
		if err != nil {
			t.Fatal(i, err)
		}
		b, err := cref.MarshalBinary()
		if err != nil {
			t.Fatal(i, err)
		}
		if !bytes.Equal(v, b) {
			fmt.Println(hex.EncodeToString(v))
			fmt.Println(hex.EncodeToString(b))

			t.Fatal(i, err)
		}

		n := NewCREF(cref.DestinationLocalReference, cref.Cause, cref.Opts)
		b2, err := n.MarshalBinary()
		if err != nil {
			t.Fatal(i, err)
		}
		if !bytes.Equal(v, b2) {
			fmt.Println(hex.EncodeToString(v))
			fmt.Println(hex.EncodeToString(b2))

			t.Fatal(i, err)
		}

		fmt.Println(cref)
	}
}
