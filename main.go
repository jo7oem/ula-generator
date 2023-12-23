package main

import (
	// sha1 is used to generate ULA prefix RFC4193 3.2.2 4).
	"crypto/sha1" //nolint:gosec
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	ULAPrefix = 0xfc // RFC4193 3.1
)

var ErrUnknown = errors.New("unknown error")

func main() {
	// NICのIPv6アドレスを取得する
	nics, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	nd := selectInterface(nics)
	if nd == nil {
		panic("no interface")
	}

	eui64, err := nd.EUI64()
	if err != nil {
		panic(err)
	}

	// 現在の時刻よりNTP Timestamp Formatのバイト列を作成する
	ntp64 := Time2Ntp64(time.Now())

	ula := genULAAddress(ntp64.Bytes(), eui64)

	fmt.Printf("%s/56\n", ula.String())
}

func genULAAddress(ntp64 []byte, eui64 []byte) net.IP {
	// NAT64とEUI-64を組み合わせてSHA1を計算する
	tmp := make([]byte, 16)
	copy(tmp, eui64)
	copy(tmp[8:], ntp64)

	dig := sha1.Sum(tmp) //nolint:gosec

	// そのSHA1の下位40bitを取り出す
	// これがULAのプレフィックスになる
	ulaPrefix := dig[12:]

	// このULAプレフィックスを使って、ULAを生成する
	ula := make([]byte, 16)
	ula[0] = ULAPrefix + 1 // global flag 0の場合は現在未定義
	copy(ula[1:], ulaPrefix)

	ip := net.IP(ula)

	return ip
}
