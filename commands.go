package noolite

// FUTURE: device presets: vid, pid, in/out endpoints etc
const (
	DeviceMTRF64USB = 1
)

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
	CmdRollColour     byte = 16
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
