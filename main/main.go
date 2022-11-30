package main

import (
	"fmt"
	"subnetcalc"
)

func main() {

	s, err := subnetcalc.Parse("10.0.0.0/16")
	if err != nil {
		panic(err)
	}

	res, err := s.AddReservation("10.0.96.0/19", "res-1")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.AddReservation("10.0.223.128/25", "res-2")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.AddReservation("192.168.1.1/22", "res-err")
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.FindFree(30)
	if res != nil {
		if err = res.Reserve("open reservation 1"); err != nil {
			fmt.Printf("Error making reservation: %v", err)
		}
	}
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)

	res, err = s.FindFree(30)
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)
	if res != nil {
		if err = res.Reserve("open reservation 2"); err != nil {
			fmt.Printf("Error making reservation: %v", err)
		}
	}

	res, err = s.FindFree(30)
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)
	if res != nil {
		if err = res.Reserve("open reservation 3"); err != nil {
			fmt.Printf("Error making reservation: %v", err)
		}
	}

	fmt.Printf("\n\n\n")

	fmt.Printf("\n\nReservations:\n")
	reservations := s.Collect(subnetcalc.FilterReserved)
	for _, r := range reservations {
		fmt.Printf(" %s %s\n", r.CIDR(), r.Reservation())
	}

	fmt.Printf("\n\n\nPlaying around with fresh data\n\n")
	s, err = subnetcalc.Parse("10.0.0.0/16")

	for i := s.Size(); i < 31; i++ {
		res, err = s.FindFree(i)
		fmt.Printf("Found free (size %d): %s %v\n", i, res.CIDR(), err)
	}

	fmt.Printf("\n\n\n")

	res, err = s.FindFree(17)
	fmt.Printf("Added reservation: %s %v\n", res.CIDR(), err)
	if res != nil {
		if err = res.Reserve("open reservation 4"); err != nil {
			fmt.Printf("Error making reservation: %v", err)
		}
	}

	fmt.Printf("\n\n\n")

	for i := s.Size(); i < 31; i++ {
		res, err = s.FindFree(i)
		fmt.Printf("Found free (size %d): %s %v\n", i, res.CIDR(), err)
	}

	fmt.Printf("\n\n\n")
}
