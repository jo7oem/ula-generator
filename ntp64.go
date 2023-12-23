package main

import (
	"encoding/binary"
	"time"
)

const UnixNtpOffset = 2208988800

type Ntp64 struct {
	Seconds  uint32
	Fraction uint32
}

func (n *Ntp64) Bytes() []byte {
	res := make([]byte, 8)
	binary.BigEndian.PutUint32(res[0:], n.Seconds)
	binary.BigEndian.PutUint32(res[4:], n.Fraction)

	return res
}

// Time2Ntp64 converts a time.Time to an NTP timestamp.
func Time2Ntp64(t time.Time) Ntp64 {
	secs := uint32(t.Unix() + UnixNtpOffset)
	frac := uint32(uint64(t.Nanosecond()) * (1 << 32) / 1e9) // 0x0001 = 約 232 ps

	return Ntp64{secs, frac}
}
