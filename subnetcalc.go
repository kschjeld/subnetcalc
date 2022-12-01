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

var ErrCouldNotParse = errors.New("could not parse subnet specification")
var ErrNotDividable = errors.New("could not divide subnet")
var ErrAlreadyReserved = errors.New("subnet is already reserved")
var ErrDidNotFindSubnet = errors.New("could not find suitable subnet")

func SelectReserved() func(s *Subnet) bool {
	return func(s *Subnet) bool {
		return s.Reservation() != ""
	}
}

func SelectAvailable() func(s *Subnet) bool {
	return func(s *Subnet) bool {
		return s.Reservation() == "" && s.subReservations == 0
	}
}

func SelectWithSize(size int) func(s *Subnet) bool {
	return func(s *Subnet) bool {
		return s.Size() == size
	}
}

func Parse(s string) (*Subnet, error) {
	c, err := toCIDR(s)
	if err != nil {
		return nil, err
	}
	subnet := &Subnet{
		cidr:   *c,
		parent: nil,
		low:    nil,
		high:   nil,
	}

	subnet.initialize()

	return subnet, nil
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
	return inetNToA(inetBToN(s.cidr.net.IP) + 1)
}

func (s *Subnet) LastIP() string {
	return inetNToA(inetSubnetLastAddress(inetBToN(s.cidr.net.IP), s.Size()) - 1)
}

func (s *Subnet) Reservation() string {
	return s.reservation
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
		_ = s.divide()
		if s.low != nil && s.low.cidr.net.Contains(cidr.net.IP) {
			return s.low.AddReservation(subnet, name)
		} else if s.high != nil {
			return s.high.AddReservation(subnet, name)
		}
	}

	return nil, ErrDidNotFindSubnet
}

func (s *Subnet) Collect(filterFunc ...func(s *Subnet) bool) []*Subnet {
	var res []*Subnet

	if s == nil {
		return res
	}

	match := true
	for _, fn := range filterFunc {
		match = fn(s)
		if !match {
			break
		}
	}

	if match {
		res = append(res, s)
	}

	res = append(res, s.low.Collect(filterFunc...)...)
	res = append(res, s.high.Collect(filterFunc...)...)

	return res
}

func (s *Subnet) FindFree(requiredSize int) (*Subnet, error) {
	if s.Size() == requiredSize && s.reservation == "" && s.subReservations == 0 {
		return s, nil
	}

	var found *Subnet
	var err error
	if s.reservation == "" && s.Size() < requiredSize {
		s.divide()

		if s.low != nil {
			found, err = s.low.FindFree(requiredSize)
			if err != nil {
				return nil, err
			}
		}

		if found == nil && s.high != nil {
			found, err = s.high.FindFree(requiredSize)
			if err != nil {
				return nil, err
			}
		}
	}

	// Top level detects nothing found, and returns error instead of nil
	if found == nil && s.parent == nil {
		return nil, ErrDidNotFindSubnet
	}

	return found, err
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

func (s *Subnet) initialize() {
	if s == nil {
		return
	}

	if s.Size() >= 31 {
		return
	}

	s.divide()
	s.low.initialize()
	s.high.initialize()
}

func (s *Subnet) divide() error {
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
		return nil, ErrCouldNotParse
	}
	return &CIDR{
		net: *snet,
	}, nil
}
