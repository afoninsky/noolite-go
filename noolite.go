package noolite

import (
	"errors"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

// Device ...
type Device struct {
	port    io.ReadWriteCloser // usb interface
	address [4]byte            // MTRF adapter address
}

// Send raw packets
func (device *Device) Send(mode, control, channel, command byte) error {
	packet := Packet{
		mode:    mode,
		control: control,
		channel: channel,
		command: command,
	}
	count, err := device.port.Write(packet.Encode())
	if err != nil {
		return err
	}
	if count != PacketLength {
		return errors.New(".Send: invalid amout of bytes written")
	}
	return nil
}

// Receive ...
func (device *Device) Receive() (Packet, error) {
	packet := Packet{}
	buf := make([]byte, PacketLength)
	if _, readError := io.ReadFull(device.port, buf); readError != nil {
		return packet, readError
	}
	if decodeError := packet.Decode(buf); decodeError != nil {
		return packet, decodeError
	}
	return packet, nil
}

// CreateDevice ...
func CreateDevice() (Device, error) {
	device := Device{}
	options := serial.OpenOptions{
		PortName:        "/dev/tty.usbserial-AL032Z5Y",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, openError := serial.Open(options)
	if openError != nil {
		return device, openError
	}
	device.port = port

	// switch to main (service) mode after device start
	if sendError := device.Send(ModeSvc, 0, 0, 0); sendError != nil {
		port.Close()
		return device, sendError
	}
	answer, receiveError := device.Receive()
	if receiveError != nil {
		return device, receiveError
	}
	if answer.mode != ModeSvc {
		return device, errors.New("device is not entered into service mode")
	}
	device.address = answer.address

	return device, nil
}
