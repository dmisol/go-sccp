package sccp

import (
	"io"
	"github.com/wmnsk/go-sccp/params"
)

func ParseOptionalParameters(b []byte) (map[uint8]*params.Optional, error) {
	p := uint8(0)
	opts := make(map[uint8]*params.Optional, 2)
	emptyOpts := func() map[uint8]*params.Optional {
		opts := make(map[uint8]*params.Optional, 0)
		return opts
	}

	for p < uint8(len(b)) {
		t := b[p]

		// No optional params
		if t == 0 {
			return opts, nil
		}

		if (p + 1) >= uint8(len(b)) {
			return emptyOpts(), io.ErrUnexpectedEOF
		}

		l := b[p+1]
		if (p + 1 + l) >= uint8(len(b)) {
			return emptyOpts(), io.ErrUnexpectedEOF
		}

		// Found option
		offset := p + 2
		o := &params.Optional {
			Tag:    t,
			Len:    l,
			Value:  b[offset:offset+l],
			Ptr:    p,
		}

		switch t {
		case params.CgPtyAddrTag:
			_, err := params.ParsePartyAddress(b[p:offset+l])
			if err != nil {
				return emptyOpts(), err
			}
		}

		opts[t] = o
		p += 2 + l
	}

	return opts, nil
}
