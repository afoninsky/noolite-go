package noolite

import (
	"reflect"
	"testing"
)

// noolite-f binded to 5th channel
// 1) enter device to service mode using remote command
// 2) bind to 5th channel
func TestNooLiteFTXBind(t *testing.T) {
	ethalonTxRemote := []byte{171, 2, 0, 0, 5, 131, 0, 1, 0, 0, 0, 0, 0, 0, 0, 54, 172}
	ethalonTxBind := []byte{171, 2, 0, 0, 10, 15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 198, 172}

	remote := &txPacket{
		mode:    modeFTx,
		control: txCtrSnd,
		channel: 5,
		command: cmdService,
	}

	bind := &txPacket{
		mode:    modeFTx,
		control: txCtrSnd,
		channel: 10,
		command: cmdBind,
	}

	if remoteEqual := reflect.DeepEqual(ethalonTxRemote, remote.Bytes()); !remoteEqual {
		t.Errorf("tx service mode packet did not match: %v", remote.Bytes())
	}

	if bindEqual := reflect.DeepEqual(ethalonTxBind, bind.Bytes()); !bindEqual {
		t.Errorf("tx bind packet did not match: %v", bind.Bytes())
	}

}

// old behaviour: manual binding
func TestNooLiteRXBind(t *testing.T) {
	ethalonTxBind := []byte{171, 1, 3, 0, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 180, 172}

	bind := &txPacket{
		mode:    modeRx,
		control: txCtrBindOn,
		channel: 5,
	}

	if bindEqual := reflect.DeepEqual(ethalonTxBind, bind.Bytes()); !bindEqual {
		t.Errorf("tx bind packet did not match: %v", bind.Bytes())
	}

}
