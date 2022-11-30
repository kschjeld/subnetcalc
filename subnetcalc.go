package subnetcalc

import (
	"errors"
	"net"
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
var ErrAlreadyReserved = errors.New("subnet is already reserved")
var ErrDidNotFindSubnet = errors.New("could not find suitable subnet")

var FilterReserved = func(s *Subnet) bool {
	return s.Reservation() != ""
}

func Parse(s string) (*Subnet, error) {
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
	return inetNToA(inetSubnetLastAddress(inetBToN(s.cidr.net.IP), s.Size()) - 1)
}

func (s *Subnet) Reservation() string {
	return s.reservation
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

func (s *Subnet) AddReservation(subnet string, name string) (*Subnet, error) {
	cidr, err := toCIDR(subnet)
	if err != nil {
		return nil, err
	}

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
			return s.low.AddReservation(subnet, name)
		} else if s.high != nil {
			return s.high.AddReservation(subnet, name)
		}
	}

	return nil, ErrDidNotFindSubnet
}

func (s *Subnet) Collect(filterFunc func(s *Subnet) bool) []*Subnet {
	var res []*Subnet

	if filterFunc(s) {
		res = append(res, s)
	}

	if s.low != nil {
		res = append(res, s.low.Collect(filterFunc)...)
	}
	if s.high != nil {
		res = append(res, s.high.Collect(filterFunc)...)
	}
	return res
}

func (s *Subnet) FindFree(requiredSize int) (*Subnet, error) {
	if s.Size() == requiredSize && s.reservation == "" && s.subReservations == 0 {
		return s, nil
	}

	if s.reservation == "" && s.Size() < requiredSize {
		s.Divide()

		if s.low != nil {
			s, err := s.low.FindFree(requiredSize)
			if s != nil || err != nil {
				return s, err
			}
		}

		if s.high != nil {
			return s.high.FindFree(requiredSize)
		}
	}

	return nil, nil
}

func (s *Subnet) Reserve(name string) error {
	if s == nil {
		return nil
	}

	if s.reservation != "" {
		if s.reservation == name {
			return nil
		}
		return ErrAlreadyReserved
	}

	s.reservation = name
	if s.parent != nil {
		s.parent.updateSubReservations()
	}

	return nil
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
			Mask: inetNToB(inetSubnetNetmask(size + 1)),
		},
	}

	high := CIDR{
		net: net.IPNet{
			IP:   inetNToB(inetSubnetLastAddress(inetBToN(s.cidr.net.IP), size+1) + 1),
			Mask: inetNToB(inetSubnetNetmask(size + 1)),
		},
	}
	return low, high, nil
}

func (s *Subnet) updateSubReservations() {
	s.subReservations = s.subReservations + 1
	if s.parent != nil {
		s.parent.updateSubReservations()
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
