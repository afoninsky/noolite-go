package noolite

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

// common packet type (tx or rx)
const (
	PacketTypeTx byte = 0
	PacketTypeRx byte = 1
)
