package subnetcalc

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type CIDR struct {
	net net.IPNet
}

type Subnet struct {
	cidr      CIDR
	parent    *Subnet
	low, high *Subnet
}

var ErrNotDividable = errors.New("unable to divide subnet")

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

func (s *Subnet) Size() int {
	size, _ := s.cidr.net.Mask.Size()
	return size
}

func (s *Subnet) FirstIP() string {
	return s.cidr.net.IP.String()
}

func (s *Subnet) LastIP() string {
	return inet_ntoa(subnet_last_address(inet_bton(s.cidr.net.IP), s.Size()) - 1)
}

func (s *Subnet) lowAndHigh() (CIDR, CIDR, error) {
	if s.low != nil && s.high != nil {
		return s.low.cidr, s.high.cidr, nil
	}

	size := s.Size()
	if size >= 32 {
		return CIDR{}, CIDR{}, ErrNotDividable
	}

	low := CIDR{
		net: net.IPNet{
			IP:   s.cidr.net.IP,
			Mask: inet_ntob(subnet_netmask(size + 1)),
		},
	}

	high := CIDR{
		net: net.IPNet{
			IP:   inet_ntob(subnet_last_address(inet_bton(s.cidr.net.IP), size+1) + 1),
			Mask: inet_ntob(subnet_netmask(size + 1)),
		},
	}
	return low, high, nil
}

func (s *Subnet) DivideRecursively(maxLevel int) {
	if s == nil {
		return
	}

	if maxLevel < 0 {
		return
	}

	s.Divide()
	s.low.DivideRecursively(maxLevel - 1)
	s.high.DivideRecursively(maxLevel - 1)
}

func (s *Subnet) Divide() error {
	if s == nil {
		return nil
	}
	if s.low != nil && s.high != nil {
		return nil
	}

	low, high, err := s.lowAndHigh()
	if err != nil {
		return err
	}

	if s.low == nil {
		s.low = &Subnet{
			cidr:   low,
			parent: s,
			low:    nil,
			high:   nil,
		}
	}

	if s.high == nil {
		s.high = &Subnet{
			cidr:   high,
			parent: s,
			low:    nil,
			high:   nil,
		}
	}
	return nil
}

func (s *Subnet) Print(onlyLeaves bool) {
	s.print(0, onlyLeaves)
}

func (s *Subnet) print(level int, onlyLeaves bool) {
	prefix := strings.Repeat(" ", level)
	if onlyLeaves && (s.low == nil && s.high == nil) || !onlyLeaves {
		fmt.Printf("%s%s (%s to %s)\n", prefix, s.cidr.net.String(), s.FirstIP(), s.LastIP())
	}

	if s.low != nil {
		s.low.print(level+1, onlyLeaves)
	}
	if s.high != nil {
		s.high.print(level+1, onlyLeaves)
	}
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
