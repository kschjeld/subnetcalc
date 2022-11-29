package main

import (
	"fmt"
	"subnetcalc"
)

func main() {

	s, err := subnetcalc.ParseSubnet("10.0.0.0/16")
	if err != nil {
		panic(err)
	}

	res, err := s.AddReservationFor("10.0.96.0/19", "res-1")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.AddReservationFor("10.0.223.128/25", "res-2")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.AddReservationFor("192.168.1.1/22", "res-err")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	fmt.Printf("\n\nRanges:\n")

	s.DivideRecursively(2)
	s.Print(true)

	res, err = s.Reserve(30, "open reservation 1")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.Reserve(30, "open reservation 2")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.Reserve(30, "open reservation 3")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	fmt.Printf("\n\nReservations:\n")
	reservations := s.Filter(subnetcalc.FilterfuncReserved)
	for _, r := range reservations {
		fmt.Printf(" %s %s\n", r.CIDR(), r.Reservation())
	}
}
