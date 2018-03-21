package noolite

const (
	txStart byte = 171
	txStop  byte = 172
)

const (
	cmdOff            byte = 0
	cmdBrightDown     byte = 1
	cmdOn             byte = 2
	cmdBrightUp       byte = 3
	cmdSwitch         byte = 4
	cmdBrightBack     byte = 5
	cmdSetBrightness  byte = 6
	cmdLoadPreset     byte = 7
	cmdSavePreset     byte = 8
	cmdUnbind         byte = 9
	cmdStopReg        byte = 10
	cmdBrightStepDown byte = 11
	cmdBrightStepUp   byte = 12
	cmdBrightReg      byte = 13
	cmdBind           byte = 15
	cmdRollColour     byte = 16
	cmdSwitchColor    byte = 17
	cmdSwitchMode     byte = 18
	cmdSpeedModeBack  byte = 19
	cmdBatteryLow     byte = 20
	cmdSendTempHumi   byte = 21
	cmdTemporaryOn    byte = 25
	cmdModes          byte = 26
	cmdReadState      byte = 128
	cmdWriteState     byte = 129
	cmdSendState      byte = 130
	cmdService        byte = 131
	cmdClearMemory    byte = 132
)

const (
	txCtrSnd      byte = 0
	txCtrSndAll   byte = 1
	txCtrRcv      byte = 2
	txCtrBindOn   byte = 3
	txCtrBindOff  byte = 4
	txCtrClear    byte = 5
	txCtrClearAll byte = 6
	txCtrUnbind   byte = 7
	txCtrSndAddr  byte = 8
)

const (
	modeTx  byte = 0
	modeRx  byte = 1
	modeFTx byte = 2
	modeFRx byte = 3
	modeSvc byte = 4
	modeUpd byte = 5
)

type txPacket struct {
	mode    byte
	control byte
	channel byte
	command byte
	address byte
}

func (tx txPacket) Bytes() []byte {
	// https://www.noo.com.by/assets/files/PDF/MTRF-64-USB.pdf
	// 2do: ability to append repeats in 3d byte
	buf := []byte{txStart, tx.mode, tx.control, 0, tx.channel, tx.command, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, txStop}

	if tx.command == cmdService {
		// allow service mode in data
		buf[7] = 1
	}

	// calculate crc
	sum := uint(0)
	for _, value := range buf[:15] {
		sum += uint(value)
	}
	buf[15] = byte(sum & 0xFF)

	return buf
}
