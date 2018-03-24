package main

import "noolite/noolite"

type cmd struct {
	value byte
	desc  string
}

// adapter modes
var modes = map[string]cmd{
	"tx":      cmd{noolite.ModeTx, "test"},
	"rx":      cmd{noolite.ModeRx, "test2"},
	"ftx":     cmd{noolite.ModeFTx, "test3"},
	"frx":     cmd{noolite.ModeFRx, "test4"},
	"service": cmd{noolite.ModeSvc, "test5"},
	"update":  cmd{noolite.ModeUpd, "test6"},
}

var txControls = map[string]cmd{
	"send":     cmd{noolite.TxCtrSnd, "desc1"},
	"sendall":  cmd{noolite.TxCtrSndAll, "desc2"},
	"receive":  cmd{noolite.TxCtrRcv, "desc3"},
	"bindon":   cmd{noolite.TxCtrBindOn, "desc4"},
	"bindoff":  cmd{noolite.TxCtrBindOff, "desc5"},
	"clear":    cmd{noolite.TxCtrClear, "desc6"},
	"clearall": cmd{noolite.TxCtrClearAll, "desc7"},
	"_unbind":  cmd{noolite.TxCtrUnbind, "desc8"},
	"sendaddr": cmd{noolite.TxCtrSndAddr, "desc9"},
}

type preset struct {
	Mode    byte
	Control byte
	Command byte
	Desc    string
}

// go run constants.go main.go tx send switch 0
// Commands ...
var Commands = map[string]preset{
	"switch-legacy": preset{noolite.ModeTx, noolite.TxCtrSnd, noolite.CmdSwitch, "Switch device state (old devices)"},
	"on-legacy":     preset{noolite.ModeTx, noolite.TxCtrSnd, noolite.CmdOn, "Enable device (old devices)"},
	"off-legacy":    preset{noolite.ModeTx, noolite.TxCtrSnd, noolite.CmdOff, "Disable device (old devices)"},
	// "off":            cmd{noolite.CmdOff, "cmd1"},
	// "brightdown":     cmd{noolite.CmdBrightDown, "cmd1"},
	// "on":             cmd{noolite.CmdOn, "cmd1"},
	// "brightup":       cmd{noolite.CmdBrightUp, "cmd1"},
	// "switch":         cmd{noolite.CmdSwitch, "cmd1"},
	// "brightback":     cmd{noolite.CmdBrightBack, "cmd1"},
	// "brightset":      cmd{noolite.CmdSetBrightness, "cmd1"},
	// "load":           cmd{noolite.CmdLoadPreset, "cmd1"},
	// "save":           cmd{noolite.CmdSavePreset, "cmd1"},
	// "unbind":         cmd{noolite.CmdUnbind, "cmd1"},
	// "stopreg":        cmd{noolite.CmdStopReg, "cmd1"},
	// "brightstepdown": cmd{noolite.CmdBrightStepDown, "cmd1"},
	// "brightstepup":   cmd{noolite.CmdBrightStepUp, "cmd1"},
	// "brightreg":      cmd{noolite.CmdBrightReg, "cmd1"},
	// "bind":           cmd{noolite.CmdBind, "cmd1"},
	// "rollcolor":      cmd{noolite.CmdRollColor, "cmd1"},
	// "switchcolor":    cmd{noolite.CmdSwitchColor, "cmd1"},
	// "switchmode":     cmd{noolite.CmdSwitchMode, "cmd1"},
	// "speedmodeback":  cmd{noolite.CmdSpeedModeBack, "cmd1"},
	// "batterylow":     cmd{noolite.CmdBatteryLow, "cmd1"},
	// "sendtemphumi":   cmd{noolite.CmdSendTempHumi, "cmd1"},
	// "tempon":         cmd{noolite.CmdTemporaryOn, "cmd1"},
	// "modes":          cmd{noolite.CmdModes, "cmd1"},
	// "readstate":      cmd{noolite.CmdReadState, "cmd1"},
	// "writestate":     cmd{noolite.CmdWriteState, "cmd1"},
	// "sendstate":      cmd{noolite.CmdSendState, "cmd1"},
	// "service":        cmd{noolite.CmdService, "cmd1"},
	// "clearmemory":    cmd{noolite.CmdClearMemory, "cmd1"},
}
