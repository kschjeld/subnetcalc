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

	reservation     string
	subReservations int
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

func (s *Subnet) CIDR() string {
	if s != nil {
		return s.cidr.net.String()
	}
	return ""
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

var ErrAlreadyReserved = errors.New("subnet is already reserved")
var ErrDidNotFindSubnet = errors.New("could not find suitable subnet")

func (s *Subnet) AddReservationFor(cidr string, name string) (*Subnet, error) {
	c, err := toCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return s.AddReservation(*c, name)
}
func (s *Subnet) AddReservation(cidr CIDR, name string) (*Subnet, error) {
	if s.cidr.net.String() == cidr.net.String() {

		if s.reservation != "" {
			if s.reservation == name {
				return s, nil
			} else {
				return s, ErrAlreadyReserved
			}
		}

		s.reservation = name
		if s.parent != nil {
			s.parent.updateSubReservations()
		}
		return s, nil
	}

	if s.cidr.net.Contains(cidr.net.IP) {
		_ = s.Divide()
		if s.low != nil && s.low.cidr.net.Contains(cidr.net.IP) {
			return s.low.AddReservation(cidr, name)
		} else if s.high != nil {
			return s.high.AddReservation(cidr, name)
		}
	}

	return nil, ErrDidNotFindSubnet
}

func (s *Subnet) updateSubReservations() {
	s.subReservations = s.subReservations + 1
	if s.parent != nil {
		s.parent.updateSubReservations()
	}
}

func (s *Subnet) Print(onlyLeaves bool) {
	s.print(0, onlyLeaves)
}

func (s *Subnet) print(level int, onlyLeaves bool) {
	prefix := strings.Repeat(" ", level)
	if onlyLeaves && (s.low == nil && s.high == nil) || !onlyLeaves {
		if onlyLeaves {
			prefix = ""
		}
		fmt.Printf("%s%s (%s to %s) %s\n", prefix, s.cidr.net.String(), s.FirstIP(), s.LastIP(), s.reservation)
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

func (s *Subnet) Filter(filterFunc func(s *Subnet) bool) []*Subnet {
	var res []*Subnet

	if filterFunc(s) {
		res = append(res, s)
	}

	if s.low != nil {
		res = append(res, s.low.Filter(filterFunc)...)
	}
	if s.high != nil {
		res = append(res, s.high.Filter(filterFunc)...)
	}
	return res
}

var FilterfuncReserved = func(s *Subnet) bool {
	return s.reservation != ""
}

func (s *Subnet) Reservation() string {
	return s.reservation
}

func (s *Subnet) Reserve(requiredSize int, name string) (*Subnet, error) {
	if s.Size() == requiredSize && s.reservation == "" && s.subReservations == 0 {
		s.reservation = name
		if s.parent != nil {
			s.parent.updateSubReservations()
		}
		return s, nil
	}

	if s.Size() < requiredSize {
		s.Divide()

		if s.low != nil {
			s, err := s.low.Reserve(requiredSize, name)
			if s != nil || err != nil {
				return s, err
			}
		}

		if s.high != nil {
			return s.high.Reserve(requiredSize, name)
		}
	}

	return nil, nil
}
