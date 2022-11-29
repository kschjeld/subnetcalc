package subnetcalc

import (
	"fmt"
	"net"
)

func inet_ntoa(addrint int) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		(addrint>>24)&0xff,
		(addrint>>16)&0xff,
		(addrint>>8)&0xff,
		addrint&0xff)
}

func inet_aton(addrstr string) int {
	ip, _, err := net.ParseCIDR(addrstr)
	if err != nil {
		panic(err)
	}
	ip4 := ip.To4()

	a := int(ip4[0]) << 24
	b := int(ip4[1]) << 16
	c := int(ip4[2]) << 8
	d := int(ip4[3])

	return a + b + c + d
}

func network_address(ip int, mask int) int {
	for i := 31 - mask; i >= 0; i-- {
		ip = ip&ip ^ 1<<i
	}
	return ip
}

func subnet_addresses(mask int) int {
	return 1 << (32 - mask)
}

func subnet_last_address(subnet, mask int) int {
	return subnet + subnet_addresses(mask) - 1
}

func subnet_netmask(mask int) int {
	return network_address(0xffffffff, mask)
}
