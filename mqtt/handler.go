package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/afoninsky/noolite-go/noolite"
	"github.com/yosssi/gmq/mqtt/client"
)

type Server struct {
	noolite *noolite.Device
	mqtt    *client.Client
}

func (s *Server) messageHandler(topicName, message []byte) {
	topicParts := strings.SplitN(string(topicName), "/", 4)
	packet := noolite.Packet{}

	if (topicParts[0] != clientID) || (len(topicParts) != 4) {
		log.Printf("ERROR: invalid command in channel - %s", topicName)
		return
	}

	// detect device type
	switch topicParts[1] {
	case modeTX:
		packet.Mode = noolite.ModeTx
	case modeRX:
		packet.Mode = noolite.ModeRx
	case modeFTX:
		packet.Mode = noolite.ModeFTx
	case modeFRX:
		packet.Mode = noolite.ModeFRx
	default:
		log.Printf("ERROR: expected noolite or noolitef device types but found - %s", topicParts[1])
		return
	}

	// get device channel
	channel, convErr := validateByteRange(topicParts[2], 1, 64)
	if convErr != nil {
		log.Printf("ERROR: invalid device channel - %s", topicParts[2])
		return
	}
	packet.Channel = channel
	if packet.Mode == noolite.ModeFTx || packet.Mode == noolite.ModeFRx {
		packet.Control = noolite.TxCtrSndAll
	} else {
		packet.Control = noolite.TxCtrSnd
	}

	// handle logic based on passed payload
	command, payload := guessCommand(message)
	switch command {
	// enters device into bind mode
	case "BIND":
		packet.Command = noolite.CmdBind

	// turns device on
	case "ON":
		packet.Command = noolite.CmdOn

	// turns device off
	case "OFF":
		packet.Command = noolite.CmdOff

	// set brightness
	case "BRIGHTNESS":
		brightness, err := validateByteRange(payload, 0, 255)
		if err != nil {
			fmt.Println(err)
			return
		}
		packet.Command = noolite.CmdSetBrightness
		packet.Data[0] = scaleBrightness(brightness)
		packet.DataFormat = 1

	// set rgb
	case "RGB":
		packet.Command = noolite.CmdSetBrightness
		packet.DataFormat = 3
		for key, value := range strings.SplitN(payload, ",", 3) {
			channel, err := validateByteRange(value, 0, 255)
			if err != nil {
				log.Printf("Error: %s", err)
				return
			}
			packet.Data[key] = channel
		}

	default:
		log.Printf("ERROR: command does not supported - %s", command)
		return
	}

	if err := s.noolite.Send(packet); err != nil {
		log.Printf("ERROR: while sending packet - %s", err)
		return
	}

	log.Printf("[SUCCESS] command %s sets on channel %v", command, channel)

	// // send the current state
	// if packet.Mode == noolite.ModeTx {
	// 	if err := s.mqtt.Publish(&client.PublishOptions{
	// 		QoS:       mqtt.QoS0,
	// 		TopicName: []byte(fmt.Sprintf(stateTopicPattern, clientID, parts[1], parts[2])),
	// 		Message:   []byte(command),
	// 	}); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}
