package subnetcalc

import (
	"fmt"
	"net"
	"strings"
)

type CIDR struct {
	net net.IPNet
}

func toCIDR(s string) (*CIDR, error) {
	_, snet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &CIDR{
		net: *snet,
	}, nil
}

type Subnet struct {
	cidr      CIDR
	parent    *Subnet
	low, high *Subnet
}

func ParseSubnet(s string) (*Subnet, error) {
	c, err := toCIDR(s)
	if err != nil {
		return nil, err
	}
	return &Subnet{
		cidr:   *c,
		parent: nil,
		low:    nil,
		high:   nil,
	}, nil
}

func (s *Subnet) Divide() {

}

func (s *Subnet) Print() {
	s.print(0)
}

func (s *Subnet) print(level int) {
	prefix := strings.Repeat(" ", level)
	fmt.Printf("%s%s\n", prefix, s.cidr.net.String())
	if s.low != nil {
		s.low.print(level + 1)
	}
	if s.high != nil {
		s.high.print(level + 1)
	}
}

func (s *Subnet) Split() error {
	return fmt.Errorf("not implemented")
}
