package noolite

import (
	"github.com/google/gousb"
)

// Noolite ..
type Noolite struct {
	ctx  *gousb.Context
	intf *gousb.Interface
	in   *gousb.InEndpoint
	out  *gousb.OutEndpoint
}

// Close ...
func (device *Noolite) Close() {
	device.ctx.Close()
	// device.intf.Close()
}

// New ...
func New(vid, pid int) (*Noolite, error) {
	noolite := Noolite{}
	vID, pID := gousb.ID(vid), gousb.ID(pid)
	noolite.ctx = gousb.NewContext()
	dev, openErr := noolite.ctx.OpenDeviceWithVIDPID(vID, pID)
	if openErr != nil {
		return nil, openErr
	}
	intf, _, openIntfErr := dev.DefaultInterface()
	if openIntfErr != nil {
		return nil, openIntfErr
	}
	noolite.intf = intf

	// MTRF-64-USB default in endpoind
	in, inErr := intf.InEndpoint(1)
	if inErr != nil {
		return nil, inErr
	}

	// MTRF-64-USB default out endpoind
	out, outErr := intf.OutEndpoint(2)
	if outErr != nil {
		return nil, outErr
	}
	noolite.in = in
	noolite.out = out

	return &noolite, nil
}
