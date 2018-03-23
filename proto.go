package noolite

import (
	"errors"
)

// PacketLength ...
const PacketLength = 17

const (
	txStart byte = 171
	txStop  byte = 172
	rxStart byte = 173
	rxStop  byte = 174
)

// Packet contains incoming decoded data OR decode data for send
type Packet struct {
	mode       byte    // <-> adapter mode
	control    byte    // <-> adapter control commands
	channel    byte    // <-> channel address
	command    byte    // <-> command
	address    [4]byte // <-> nooliteF device address
	repeat     byte    // -> amount of repeats + 2 {not realized}
	dataFormat byte    // <- format of incoming data
	data       [4]byte // <- incoming data
	toggle     byte    // <- nooliteF-tx - amount of packets to receive, nooliteF-rx/noolite - unique command id
}

func crc(input []byte) byte {
	sum := uint(0)
	for _, value := range input {
		sum += uint(value)
	}
	return byte(sum & 0xFF)
}

// Encode creates rx packet for MTRF device
func (p Packet) Encode() []byte {
	buf := []byte{txStart, 0, p.control, 0, p.channel, p.command, 0, p.address[0], p.address[1], p.address[2], p.address[3], 0, 0, 0, 0, 0, txStop}

	// add mode + repeats flag
	buf[1] = (p.repeat << 6) | p.mode

	// enable service mode flag in case of service command
	if p.command == CmdService {
		buf[7] = 1
	}
	buf[15] = crc(buf[:15])
	return buf
}

// Decode parses incoming data from MTRF adapter into readable structure
func (p *Packet) Decode(buf []byte) error {
	if len(buf) != PacketLength || buf[0] != rxStart || buf[PacketLength-1] != rxStop {
		return errors.New(".Decode: invalid packet format")
	}
	if crc(buf[:15]) != buf[15] {
		return errors.New(".Decode: invalid crc")
	}

	p.mode = buf[1]
	p.control = buf[2]
	p.toggle = buf[3]
	p.channel = buf[4]
	p.command = buf[5]
	p.dataFormat = buf[6]
	p.data = [4]byte{buf[7], buf[8], buf[9], buf[10]}
	p.address = [4]byte{buf[11], buf[12], buf[13], buf[14]}

	return nil
}
