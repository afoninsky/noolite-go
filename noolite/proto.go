package noolite

import (
	"errors"
)

// PacketLength ...
const PacketLength = 17

const (
	txStart byte = 171 // start byte for outgoing packet
	txStop  byte = 172 // stop byte for outgoing packet
	rxStart byte = 173 // start byte for incoming packet
	rxStop  byte = 174 // stop byte for incoming packet
)

// Packet contains incoming decoded data OR decode data for send
type Packet struct {
	Mode       byte    // <-> adapter mode
	Control    byte    // <-> adapter control commands
	Channel    byte    // <-> channel address
	Command    byte    // <-> command
	Address    [4]byte // <-> nooliteF device address
	Repeat     byte    // -> amount of repeats + 2 {not realized}
	DataFormat byte    // <-> format of incoming data or amout of outgoing
	Data       [4]byte // <- incoming data
	Toggle     byte    // <- nooliteF-tx - amount of packets to receive,nooliteF-rx/noolite - unique command id
	Type       byte
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
	buf := []byte{txStart, 0, p.Control, 0, p.Channel, p.Command, p.DataFormat, p.Data[0], p.Data[1], p.Data[2], p.Data[3], p.Address[0], p.Address[1], p.Address[2], p.Address[3], 0, txStop}

	// add mode + repeats flag
	buf[1] = (p.Repeat << 6) | p.Mode

	// enable service mode flag in case of service command
	if p.Command == CmdService {
		buf[7] = 1
	}

	// count crc
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

	p.Mode = buf[1]
	p.Control = buf[2]
	p.Toggle = buf[3]
	p.Channel = buf[4]
	p.Command = buf[5]
	p.DataFormat = buf[6]
	p.Data = [4]byte{buf[7], buf[8], buf[9], buf[10]}
	p.Address = [4]byte{buf[11], buf[12], buf[13], buf[14]}

	switch p.Mode {
	case ModeTx:
		p.Type = PacketTypeTx
	case ModeRx:
		p.Type = PacketTypeTx
	case ModeFTx:
		p.Type = PacketTypeRx
	case ModeFRx:
		p.Type = PacketTypeRx
	default:
		return errors.New(".Decode: unsupported packet type")
	}
	return nil
}
