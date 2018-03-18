package main

import (
	"log"

	"github.com/google/gousb"
)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	vid, pid := gousb.ID(0x0403), gousb.ID(0x6001)
	dev, err := ctx.OpenDeviceWithVIDPID(vid, pid)
	if err != nil {
		log.Fatalf("OpenDeviceWithVIDPID(): %v", err)
	}
	intf, done, err := dev.DefaultInterface()
	if err != nil {
		log.Fatalf("DefaultInterface(): %v", err)
	}
	defer done()

	// MTRF-64-USB default in endpoind
	epIn, err := intf.InEndpoint(1)
	if err != nil {
		log.Fatalf("%s.InEndpoint(1): %v", intf, err)
	}

	// MTRF-64-USB default out endpoind
	epOut, err := intf.OutEndpoint(2)
	if err != nil {
		log.Fatalf("%s.OutEndpoint(2): %v", intf, err)
	}

}
