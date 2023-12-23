package main

import (
	"testing"
	"time"
)

func TestTime2Ntp64(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name  string
		Param time.Time
		Want  Ntp64
	}{
		{
			Name:  "Zero",
			Param: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			Want:  Ntp64{Seconds: 2208988800, Fraction: 0},
		},
		{
			Name:  "+1sec",
			Param: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(1 * time.Second),
			Want:  Ntp64{Seconds: 2208988800 + 1, Fraction: 0},
		},
		{
			Name:  "1ns",
			Param: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(1 * time.Nanosecond),
			Want:  Ntp64{Seconds: 2208988800, Fraction: 0x00000004},
		},
		{
			Name:  "500ms",
			Param: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(500 * time.Millisecond),
			Want:  Ntp64{Seconds: 2208988800, Fraction: 0x80000000},
		},
		{
			Name:  "0.999 999 999 sec",
			Param: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Add(1*time.Second - 1*time.Nanosecond),
			Want:  Ntp64{Seconds: 2208988800, Fraction: 0xFFFFFFFB},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			got := Time2Ntp64(tc.Param)
			if got != tc.Want {
				t.Errorf("got %v, want %v", got, tc.Want)
				t.Errorf("\n got\t %08b\nwant\t %08b", got.Bytes(), tc.Want.Bytes())
			}
		})
	}
}
