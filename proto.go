package noolite

import (
	"errors"
)

const packetLength = 17

const (
	txStart byte = 171
	txStop  byte = 172
	rxStart byte = 173
	rxStop  byte = 174
)

// commands to execute
const (
	CmdOff            byte = 0
	CmdBrightDown     byte = 1
	CmdOn             byte = 2
	CmdBrightUp       byte = 3
	CmdSwitch         byte = 4
	CmdBrightBack     byte = 5
	CmdSetBrightness  byte = 6
	CmdLoadPreset     byte = 7
	CmdSavePreset     byte = 8
	CmdUnbind         byte = 9
	CmdStopReg        byte = 10
	CmdBrightStepDown byte = 11
	CmdBrightStepUp   byte = 12
	CmdBrightReg      byte = 13
	CmdBind           byte = 15
	CmdRollColor      byte = 16
	CmdSwitchColor    byte = 17
	CmdSwitchMode     byte = 18
	CmdSpeedModeBack  byte = 19
	CmdBatteryLow     byte = 20
	CmdSendTempHumi   byte = 21
	CmdTemporaryOn    byte = 25
	CmdModes          byte = 26
	CmdReadState      byte = 128
	CmdWriteState     byte = 129
	CmdSendState      byte = 130
	CmdService        byte = 131
	CmdClearMemory    byte = 132
)

// controller state
const (
	TxCtrSnd      byte = 0
	TxCtrSndAll   byte = 1
	TxCtrRcv      byte = 2
	TxCtrBindOn   byte = 3
	TxCtrBindOff  byte = 4
	TxCtrClear    byte = 5
	TxCtrClearAll byte = 6
	TxCtrUnbind   byte = 7
	TxCtrSndAddr  byte = 8

	RxCtrCmdSuccess    byte = 0
	RxCtrTimeout       byte = 1
	RxCtrError         byte = 2
	RxCtrBindCommplete byte = 3
)

// controller mode
const (
	ModeTx  byte = 0
	ModeRx  byte = 1
	ModeFTx byte = 2
	ModeFRx byte = 3
	ModeSvc byte = 4
	ModeUpd byte = 5
)

// Packet contains incoming decoded data OR decode data for send
type Packet struct {
	mode    byte    // <-> adapter mode
	control byte    // <-> adapter control commands
	channel byte    // <-> channel address
	command byte    // <-> command
	address [4]byte // <-> nooliteF device address
	// repeat  byte    // -> amount of repeats + 2 {not realized}
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
	buf := []byte{txStart, p.mode, p.control, 0, p.channel, p.command, 0, p.address[0], p.address[1], p.address[2], p.address[3], 0, 0, 0, 0, 0, txStop}

	// 2do: set repeats
	// buf[2] |= (tx.repeat >> 6)

	// enable service mode flag in case of service command
	if p.command == CmdService {
		buf[7] = 1
	}
	buf[15] = crc(buf[:15])
	return buf
}

// Decode parses incoming data from MTRF adapter into readable structure
func (p *Packet) Decode(buf []byte) error {
	if len(buf) != packetLength || buf[0] != rxStart || buf[packetLength-1] != rxStop {
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
