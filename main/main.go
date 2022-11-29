package main

import (
	"fmt"
	"net"
	"subnetcalc"
)

func main() {

	s, err := subnetcalc.ParseSubnet("10.0.0.0/16")
	if err != nil {
		panic(err)
	}
	//s.Print()

	//ps := "10.0.0.0/16"
	//pnm := 16
	//print(ps, pnm)
	//
	//// low
	//print(ps, pnm+1)
	//
	//// high
	//pip := inet_aton(ps)
	//lip := subnet_last_address(pip, pnm+1)
	//fip := lip + 1
	//hip := fmt.Sprintf("%s/%d", inet_ntoa(fip), pnm+1)
	//print(hip, pnm+1)

	s.DivideRecursively(2)
	s.Print(true)
}

func print(s string, sm int) {
	fmt.Printf("\n\n* %s *\n", s)
	fmt.Printf("IP: %s\n", inet_ntoa(inet_aton(s)))
	fmt.Printf("Mask: %s\n", inet_ntoa(subnet_netmask(sm)))
	fmt.Printf("Last address: %s\n", inet_ntoa(subnet_last_address(inet_aton(s), sm)))
}

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
