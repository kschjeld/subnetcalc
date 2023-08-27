package main

import (
	"fmt"
	"github.com/kschjeld/subnetcalc"
)

func main() {

	// Parse network CIDR
	nw, err := subnetcalc.Parse("10.0.0.0/16")
	if err != nil {
		panic(err)
	}

	// Register a pre-defined subnet
	_, err = nw.AddReservation("10.0.0.0/28", "Preallocated 28")
	if err != nil {
		panic(fmt.Errorf("error adding predefined subnet: %v", err))
	}
	_, err = nw.AddReservation("10.0.0.0/28", "subnet-2")
	if err != nil {
		// This was supposed to fail, it was already defined
	} else {
		panic(fmt.Errorf("did not fail adding predefined subnet: %v", err))
	}

	// Find a new 24-subnet
	s, err := nw.FindFree(24)
	if err != nil {
		panic(fmt.Errorf("error finding free subnet: %v", err))
	}

	// And reserve it
	err = s.Reserve("My new 24")
	if err != nil {
		panic(fmt.Errorf("error reserving subnet: %v", err))
	}

	// Try to reserve it again with another reservation name - should fail
	err = s.Reserve("My stolen new 24") // Using same reservation name is not supposed to fail
	if err == nil {
		panic(fmt.Errorf("did not fail reserving subnet: %v", err))
	}

	// Find another 28, to see that we pack them good and tight
	if s, err = nw.FindFreeAndReserve(28, "Another 28"); err != nil {
		panic(fmt.Errorf("error finding free 28 subnet: %v", err))
	}

	// List and print the reservations of our subnet
	fmt.Printf("The network %s has following subnets/reservations:\n", nw.CIDR())
	for _, s := range nw.Collect(subnetcalc.SelectReserved()) {
		fmt.Printf(" - %s (reservation name: %s)\n", s.CIDR(), s.Reservation())
	}
}
