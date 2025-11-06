package sccp

import (
	"fmt"
	"testing"
)

var mockUDT = []byte{9, 0, 12, 2, 6, 4, 67, 178, 58, 254, 3, 0, 1, 49, 4, 67, 240, 60, 254}

func TestUDT(t *testing.T) {
	u, err := ParseUDT(mockUDT)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u.Data)
}
