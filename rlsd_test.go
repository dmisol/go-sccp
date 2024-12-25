package sccp

import (
	"log"
	"testing"
)

func TestRlsCause(t *testing.T) {
	rlsd := new(RLSD)
	rlsd.RlsCause = 0b00001101
	if rlsd.RlsCauseName() != "Expiration of receive inactivity timer" {
		log.Fatal("Rls cause name was not received")
	}
}