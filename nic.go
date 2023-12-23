package main

import (
	"cmp"
	"fmt"
	"net"
	"slices"
)

const (
	hasMACAddrFlag       = 0x01
	hasLinkLocalAddrFlag = 0x02
	hasULAFlag           = 0x04
	hasGlobalAddrFlag    = 0x08
)

type NICScore struct {
	Flags uint8
	NIC   *net.Interface
}

func (ns *NICScore) EUI64() ([]byte, error) {
	eui64 := make([]byte, 8)

	switch {
	case ns.Flags&hasLinkLocalAddrFlag != 0:
		// ここでエラーになることはない
		addrs, err := ns.NIC.Addrs()
		if err != nil {
			return nil, fmt.Errorf("failed to get %s addresses: %w", ns.NIC.Name, err)
		}

		for _, addr := range addrs {
			ip, ok := ToIPv6Addr(addr)
			if !ok {
				continue
			}

			if ip.IsLinkLocalUnicast() {
				copy(eui64, ip[8:])

				break
			}
		}
	case ns.Flags&hasMACAddrFlag != 0:
		mac := []byte(ns.NIC.HardwareAddr)

		if len(mac) == 8 {
			return mac, nil
		}

		return GenEUI64(mac), nil

	default:
		return nil, ErrUnknown
	}

	return eui64, nil
}

func GenEUI64(mac []byte) []byte {
	eui64 := make([]byte, 8)
	mac[0] ^= 0x02

	copy(eui64, mac[:3])
	eui64[3] = 0xFF
	eui64[4] = 0xFE
	copy(eui64[5:], mac[3:])

	return eui64
}

func ToIPv6Addr(adr net.Addr) (net.IP, bool) {
	ip, ok := adr.(*net.IPNet)
	if !ok {
		return nil, false
	}

	v4 := ip.IP.To4()
	v6 := ip.IP.To16()

	if v4 != nil || v6 == nil {
		return nil, false
	}

	return v6, true
}

func isULA(ip net.IP) bool {
	if len(ip) != 16 {
		return false
	}

	return ip[0]&0xE == ULAPrefix
}

func selectInterface(nics []net.Interface) *NICScore {
	filtered := make([]*NICScore, 0, len(nics))

	for _, nic := range nics {
		nic := nic

		nd, err := genNICDetail(&nic)
		if err != nil || nd == nil {
			continue
		}

		filtered = append(filtered, nd)
	}

	slices.SortStableFunc(filtered, func(i, j *NICScore) int {
		// 優先度の高いNICが先頭にきてほしい
		return cmp.Compare(j.Flags, i.Flags)
	})

	if filtered[0].Flags == 0 {
		return nil
	}

	return filtered[0]
}

func genNICDetail(nic *net.Interface) (*NICScore, error) {
	adders, err := nic.Addrs()
	if err != nil {
		return nil, fmt.Errorf("failed to get %s addresses: %w", nic.Name, err)
	}

	nd := NICScore{
		Flags: 0,
		NIC:   nic,
	}

	if len(nic.HardwareAddr) == 6 || len(nic.HardwareAddr) == 8 {
		nd.Flags |= hasMACAddrFlag
	}

	for _, addr := range adders {
		ip, ok := ToIPv6Addr(addr)
		if !ok {
			continue
		}

		if isULA(ip) {
			nd.Flags |= hasULAFlag

			// ULAとグローバルアドレスは区別する
			continue
		}

		if ip.IsGlobalUnicast() {
			nd.Flags |= hasGlobalAddrFlag

			continue
		}

		if ip.IsLinkLocalUnicast() {
			nd.Flags |= hasLinkLocalAddrFlag

			continue
		}
	}

	return &nd, nil
}
