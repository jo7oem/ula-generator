package main

import (
	"bytes"
	"testing"
)

func TestGenEUI64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Param []byte
		Want  []byte
	}{
		{
			Name:  "case1",
			Param: []byte{0x00, 0x12, 0x34, 0xAB, 0xCD, 0xEF},
			Want:  []byte{0x02, 0x12, 0x34, 0xFF, 0xFE, 0xAB, 0xCD, 0xEF},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			got := GenEUI64(tc.Param)
			if !bytes.Equal(got, tc.Want) {
				t.Errorf("got %v, want %v", got, tc.Want)
			}
		})
	}
}
