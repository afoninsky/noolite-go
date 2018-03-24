package noolite

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("%s != %s", a, b)
	}
}

// noolite-f binded to 5th channel
// - enter device to service mode using remote command
// - bind to 5th channel
// - wait for success response
func TestBindRemote(t *testing.T) {
	ethalonTxRemote := []byte{171, 2, 0, 0, 5, 131, 0, 1, 0, 0, 0, 0, 0, 0, 0, 54, 172}
	ethalonTxBind := []byte{171, 2, 0, 0, 10, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 198, 172}
	ethalonRxBind := []byte{173, 2, 3, 0, 5, 130, 0, 1, 1, 1, 1, 2, 2, 2, 2, 0, 174}
	ethalonRxBind[15] = crc(ethalonRxBind[:15])

	remote := &Packet{
		mode:    ModeFTx,
		control: TxCtrSnd,
		channel: 5,
		command: CmdService,
	}

	bind := &Packet{
		mode:    ModeFTx,
		control: TxCtrSnd,
		channel: 10,
		command: CmdBind,
	}

	receive := &Packet{}

	if remoteEqual := reflect.DeepEqual(ethalonTxRemote, remote.Encode()); !remoteEqual {
		t.Errorf("tx service mode packet did not match: %v", remote.Encode())
	}

	if bindEqual := reflect.DeepEqual(ethalonTxBind, bind.Encode()); !bindEqual {
		t.Errorf("tx bind packet did not match: %v", bind.Encode())
	}

	if rxErr := receive.Decode(ethalonRxBind); rxErr != nil {
		t.Error(rxErr)
	}
	assertEqual(t, receive.mode, ModeFTx)
	assertEqual(t, receive.control, RxCtrBindCommplete)
	assertEqual(t, receive.toggle, byte(0))
	assertEqual(t, receive.channel, byte(5))
	assertEqual(t, receive.command, CmdSendState)
	assertEqual(t, receive.dataFormat, byte(0))

}

// old behaviour: manual binding
// - send command
// - wait for success command in 40 seconds
func TestBindLocal(t *testing.T) {
	ethalonTxBind := []byte{171, 1, 3, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 180, 172}

	bind := &Packet{
		mode:    ModeRx,
		control: TxCtrBindOn,
		channel: 5,
	}

	if bindEqual := reflect.DeepEqual(ethalonTxBind, bind.Encode()); !bindEqual {
		t.Errorf("tx bind packet did not match: %v", bind.Encode())
	}

}

// // completly unbind all channels from adapter
func TestClearChannels(t *testing.T) {
	receiver := [4]byte{170, 85, 170, 85}
	ethalonClearNoolite := []byte{171, 1, 6, 0, 0, 0, 0, 170, 85, 170, 85, 0, 0, 0, 0, 176, 172}
	ethalonClearNooliteF := []byte{171, 3, 6, 0, 0, 0, 0, 170, 85, 170, 85, 0, 0, 0, 0, 178, 172}

	clear := &Packet{
		mode:    ModeRx,
		control: TxCtrClearAll,
		address: receiver,
	}

	clearF := &Packet{
		mode:    ModeFRx,
		control: TxCtrClearAll,
		address: receiver,
	}

	if clearEqual := reflect.DeepEqual(ethalonClearNoolite, clear.Encode()); !clearEqual {
		t.Errorf("clear packet did not match: %v", clear.Encode())
	}

	if clearFEqual := reflect.DeepEqual(ethalonClearNooliteF, clearF.Encode()); !clearFEqual {
		t.Errorf("clearF packet did not match: %v", clearF.Encode())
	}

}

func TestSendOnOff(t *testing.T) {
	ethalonOnSingle := []byte{171, 2, 0, 0, 10, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 185, 172}
	ethalonOffBroadcast := []byte{171, 2, 1, 0, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 184, 172}

	onSingle := &Packet{
		mode:    ModeFTx,
		control: TxCtrSnd,
		command: CmdOn,
		channel: 10,
	}

	offBroadcast := &Packet{
		mode:    ModeFTx,
		control: TxCtrSndAll,
		command: CmdOff,
		channel: 10,
	}

	if onSingleEqual := reflect.DeepEqual(ethalonOnSingle, onSingle.Encode()); !onSingleEqual {
		t.Errorf("onSingle packet did not match: %v", onSingle.Encode())
	}

	if offBroadcastEqual := reflect.DeepEqual(ethalonOffBroadcast, offBroadcast.Encode()); !offBroadcastEqual {
		t.Errorf("offBroadcast packet did not match: %v", offBroadcast.Encode())
	}

}
