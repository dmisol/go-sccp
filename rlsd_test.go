package sccp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

var mockRLSD = []byte{0x4, 0x0, 0x3, 0x70, 0x0, 0x80, 0x2a, 0x0, 0x0}

func TestRLSD(t *testing.T) {
	r, err := ParseRLSD(mockRLSD)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(mockRLSD, b) {
		fmt.Println(hex.EncodeToString(mockRLSD))
		fmt.Println(hex.EncodeToString(b))

		t.Fatal(err)
	}

	n := NewRLSD(r.DestinationLocalReference, r.SourceLocalReference, r.Cause, r.Opts)
	b2, err := n.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(mockRLSD, b2) {
		fmt.Println(hex.EncodeToString(mockRLSD))
		fmt.Println(hex.EncodeToString(b2))

		t.Fatal(err)
	}

	fmt.Println(r)
}
