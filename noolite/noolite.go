package noolite

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"

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
		MinimumReadSize: 1,
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
	buf := packet.Encode()
	fmt.Println("->", buf)
	count, err := device.Port.Write(buf)
	if err != nil {
		return err
	}
	if count != PacketLength {
		return errors.New(".Send: invalid amout of bytes written")
	}
	return nil
}

func (device *Device) Listen(handler func(message Packet)) {
	// read by one symbol and compose result into packet
	reader := bufio.NewReader(device.Port)
	accumulator := []byte{}
	// func (b *Buffer) Write(p []byte) (n int, err error)
	for {
		buf, readError := reader.ReadBytes(rxStop)
		if readError != nil {
			log.Printf("ERROR: receive failed - %s", readError)
			accumulator = []byte{}
			continue
		}
		fmt.Println("<-", buf)
		for _, item := range buf {
			accumulator = append(accumulator, item)
		}
		if len(accumulator) < PacketLength {
			// stop-byte found in the middle of a packet, ex.: checksum
			continue
		}
		packet := Packet{}
		if decodeError := packet.Decode(accumulator); decodeError != nil {
			log.Printf("ERROR: decode failed - %s", decodeError)
			accumulator = []byte{}
			continue
		}
		go handler(packet)
		accumulator = []byte{}
	}
}

// Receive ...
// func (device *Device) Receive() (Packet, error) {
// 	fmt.Println("check something...")
// 	buf := make([]byte, 33)
// 	if _, err := io.ReadAtLeast(device.Port, buf, 4); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("got something: %s\n", buf)

// 	packet := Packet{}
// 	// buf := make([]byte, PacketLength)
// 	// if _, readError := io.ReadFull(device.Port, buf); readError != nil {
// 	// 	return packet, readError
// 	// }
// 	// if decodeError := packet.Decode(buf); decodeError != nil {
// 	// 	return packet, decodeError
// 	// }
// 	return packet, nil
// }

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
