package noolite

import (
	"errors"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

// Device implements MTRF structure
type Device struct {
	Port    io.ReadWriteCloser // usb interface
	Address [4]byte            // MTRF adapter address
	Mode    byte               // current device mode
}

// CreateDevice returns MTRF device
func CreateDevice(portName string) (Device, error) {

	device := Device{}
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}
	port, openError := serial.Open(options)
	if openError != nil {
		return device, openError
	}
	device.Port = port

	// switch to main (service) mode after device start
	// if modeErr := device.enterServiceMode(); modeErr != nil {
	// 	port.Close()
	// 	return device, modeErr
	// }

	return device, nil
}

// Send raw packets
func (device *Device) Send(packet Packet) error {
	count, err := device.Port.Write(packet.Encode())
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
	if _, readError := io.ReadFull(device.Port, buf); readError != nil {
		return packet, readError
	}
	if decodeError := packet.Decode(buf); decodeError != nil {
		return packet, decodeError
	}
	return packet, nil
}

// func (device *Device) enterServiceMode() error {
// 	if sendError := device.Send(ModeSvc, 0, 0, 0); sendError != nil {
// 		return sendError
// 	}
// 	answer, receiveError := device.Receive()
// 	if receiveError != nil {
// 		return receiveError
// 	}
// 	if answer.Mode != ModeSvc {
// 		return errors.New("device is not entered into service state")
// 	}
// 	device.Address = answer.Address
// 	return nil
// }

// Close ...
func (device *Device) Close() {
	device.Port.Close()
}
