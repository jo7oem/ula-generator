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

	ULAPrefix = 0xfc // RFC4193 3.1
)

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

	//        algorithm.  If an EUI-64 does not exist, one can be created from
	//        a 48-bit MAC address as specified in [ADDARCH].  If an EUI-64
	//        cannot be obtained or created, a suitably unique identifier,
	//        local to the node, should be used (e.g., system serial number).
	fmt.Printf("NIC: %s\n", nd.NIC.Name)
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

type NICScore struct {
	Flags uint8
	NIC   *net.Interface
}

func selectInterface(NICs []net.Interface) *NICScore {
	filtered := make([]*NICScore, 0, len(NICs))

	for _, nic := range NICs {
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

	if nic.HardwareAddr != nil {
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
