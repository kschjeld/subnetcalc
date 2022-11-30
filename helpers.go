package subnetcalc

import (
	"fmt"
)

func inetNToA(addrint int) string {
	return fmt.Sprintf("%d.%d.%d.%d", (addrint>>24)&0xff, (addrint>>16)&0xff, (addrint>>8)&0xff, addrint&0xff)
}

func inetNToB(addrint int) []byte {
	return []byte{
		byte((addrint >> 24) & 0xff),
		byte((addrint >> 16) & 0xff),
		byte((addrint >> 8) & 0xff),
		byte(addrint & 0xff),
	}
}

func inetBToN(ip4 []byte) int {
	a := int(ip4[0]) << 24
	b := int(ip4[1]) << 16
	c := int(ip4[2]) << 8
	d := int(ip4[3])

	return a + b + c + d
}

func inetNetworkAddress(ip int, mask int) int {
	for i := 31 - mask; i >= 0; i-- {
		ip = ip&ip ^ 1<<i
	}
	return ip
}

func inetSubnetAddresses(mask int) int {
	return 1 << (32 - mask)
}

func inetSubnetLastAddress(subnet, mask int) int {
	return subnet + inetSubnetAddresses(mask) - 1
}

func inetSubnetNetmask(mask int) int {
	return inetNetworkAddress(0xffffffff, mask)
}
