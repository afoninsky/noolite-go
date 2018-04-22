package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/afoninsky/noolite-go/noolite"
	"github.com/yosssi/gmq/mqtt/client"
)

type Server struct {
	noolite *noolite.Device
	mqtt    *client.Client
}

func (s *Server) messageHandler(topicName, message []byte) {
	parts := strings.Split(string(topicName), "/")
	packet := noolite.Packet{}

	// detect device type
	switch parts[1] {
	case nooliteID:
		packet.Mode = noolite.ModeTx
	case nooliteFID:
		packet.Mode = noolite.ModeFTx
	default:
		fmt.Println("expected noolite or noolitef device types but found:", parts[2])
		return
	}

	// get device channel
	channel, convErr := strconv.Atoi(parts[2])
	if convErr != nil || channel < 0 || channel > maxChannelsLength-1 {
		fmt.Println("invalid device channel:", parts[2])
		return
	}
	packet.Channel = byte(channel)

	command := strings.ToUpper(fmt.Sprintf("%s", message))
	switch command {

	case "BIND":

		if packet.Mode == noolite.ModeTx {
			// mosquitto_pub -t "home/noolite/0/set" -m bind
			// enters device into bind mode (old mode)
			bind := packet
			bind.Control = noolite.TxCtrSnd
			bind.Command = noolite.TxCtrBindOn
			if err := s.noolite.Send(bind); err != nil {
				fmt.Println(err)
				return
			}

		} else {
			// set command "bind remotely"
			remote := packet
			remote.Command = noolite.CmdService

			// perform binding
			bind := packet
			bind.Command = noolite.CmdBind
			// 2do: send both commands
		}

	case "ON": // turns device on
		on := packet
		on.Command = noolite.CmdOn
		if err := s.noolite.Send(on); err != nil {
			fmt.Println(err)
			return
		}
	case "OFF": // turns device off
		off := packet
		off.Command = noolite.CmdOff
		if err := s.noolite.Send(off); err != nil {
			fmt.Println(err)
			return
		}
	default:
		// check if its brightness
		// brightness, err := strconv.Atoi(command)
		// 2do: brightness mode, rgb mode

		// if err == nil && brightness >= 0 && brightness <= 255 {
		// 	bri := packet
		// 	// set brightness
		// 	bri.Command = noolite.CmdSetBrightness
		// 	// 2do 0..155
		// 	if err := s.noolite.Send(bri); err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// 	return
		// }
		// inform command does not supported
		fmt.Println("command does not supported:", command)
		return
	}

	fmt.Println("command", command, "sets on channel", channel)

	// // send the current state
	// if packet.Mode == noolite.ModeTx {
	// 	if err := s.mqtt.Publish(&client.PublishOptions{
	// 		QoS:       mqtt.QoS0,
	// 		TopicName: []byte(fmt.Sprintf(stateTopicPattern, parts[1], parts[2])),
	// 		Message:   []byte(command),
	// 	}); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}
