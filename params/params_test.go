// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params_test

import (
	"encoding"
	"io"
	"testing"

	"github.com/dmisol/go-sccp/params"
	"github.com/pascaldekloe/goe/verify"
)

type serializable interface {
	encoding.BinaryMarshaler
	MarshalLen() int
}

type decodeFunc func([]byte) (serializable, error)

var testcases = []struct {
	description string
	structured  serializable
	serialized  []byte
	decodeFunc
}{
	{
		description: "PartyAddress",
		structured: params.NewPartyAddress(
			0x12, 0, 6, 0, // GTI, SPC, SSN, TT
			1, 1, 4, // NP, ES, NAI
			[]byte{
				0x21, 0x43, 0x65, 0x87, 0x09,
			},
		),
		serialized: []byte{
			0x0a, 0x12, 0x06, 0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09,
		},
		decodeFunc: func(b []byte) (serializable, error) {
			v, err := params.ParsePartyAddress(b)
			if err != nil {
				return nil, err
			}

			return v, nil
		},
	}, {
		description: "PartyAddress/2-bytes",
		structured: params.NewPartyAddress(
			0x42, 0, 6, 0x00, // Indicator, SPC, SSN, TT
			0x00, 0x00, 0x00, // NP, ES, NAI
			nil, // GlobalTitleInformation
		),
		serialized: []byte{
			0x02, 0x42, 0x06,
		},
		decodeFunc: func(b []byte) (serializable, error) {
			v, err := params.ParsePartyAddress(b)
			if err != nil {
				return nil, err
			}

			return v, nil
		},
	},
}

func TestStructuredParams(t *testing.T) {
	t.Helper()

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				prm, err := c.decodeFunc(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := prm, c.structured; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}

				if got, want := b, c.serialized; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("Len", func(t *testing.T) {
				if got, want := c.structured.MarshalLen(), len(c.serialized); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})
		})
	}
}

func TestPartialStructuredParams(t *testing.T) {
	for _, c := range testcases {
		for i := range c.serialized {
			partial := c.serialized[:i]
			_, err := c.decodeFunc(partial)
			if err != io.ErrUnexpectedEOF {
				t.Errorf("%#x: got error %v, want unexpected EOF", partial, err)
			}
		}
	}
}
