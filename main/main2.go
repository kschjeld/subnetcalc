package main

//package main
//
//import (
//	"fmt"
//	"net"
//	"subnetcalc"
//)
//
//func main() {
//
//	_, err := subnetcalc.ParseSubnet("10.0.0.0/16")
//	if err != nil {
//		panic(err)
//	}
//	//s.Print()
//
//	i := inet_aton("10.0.0.0/16")
//	//m := inet_aton("255.255.0.0/16")
//	m := subnet_netmask(uint32(16))
//
//	fmt.Printf("IP: %s\n", inet_ntoa(i))
//	fmt.Printf("Mask: %s\n", inet_ntoa(m))
//
//}
//
//func inet_ntoa(addrint uint32) string {
//	return fmt.Sprintf("%d.%d.%d.%d",
//		((addrint >> 24) & 0xff),
//		((addrint >> 16) & 0xff),
//		((addrint >> 8) & 0xff),
//		(addrint & 0xff))
//}
//
//func inet_aton(addrstr string) uint32 {
//	ip, _, err := net.ParseCIDR(addrstr)
//	if err != nil {
//		panic(err)
//	}
//	ip4 := ip.To4()
//
//	a := uint32(ip4[0]) << 24
//	b := uint32(ip4[1]) << 16
//	c := uint32(ip4[2]) << 8
//	d := uint32(ip4[3])
//
//	return a + b + c + d
//}
//
//func network_address(ip uint32, mask uint32) uint32 {
//	for i := 31 - mask; i >= 0; i-- {
//		fmt.Printf("i=%d\n", i)
//		ip = ip&ip ^ 1<<i
//	}
//	return ip
//}
//
//func subnet_addresses(mask uint32) uint32 {
//	return 1 << (32 - mask)
//}
//
//func subnet_last_address(subnet, mask uint32) uint32 {
//	return subnet + subnet_addresses(mask) - 1
//}
//
//func subnet_netmask(mask uint32) uint32 {
//	return network_address(0xffffffff, mask)
//}
