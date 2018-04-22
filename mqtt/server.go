package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/afoninsky/noolite-go/noolite"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// https://www.home-assistant.io/components/light.mqtt/

func main() {

	willTopic := "home/noolite/status"

	// handle interrupt signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// open noolite connected device
	nooDevice, nooErr := noolite.CreateDevice()
	if nooErr != nil {
		fmt.Println(nooErr)
		os.Exit(1)
	}
	defer nooDevice.Close()

	// connect to MQTT server
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
			os.Exit(1)
		},
	})
	defer cli.Terminate()
	cliErr := cli.Connect(&client.ConnectOptions{
		Network:     "tcp",
		Address:     "localhost:1883",
		ClientID:    []byte("example-client"),
		WillTopic:   []byte(willTopic),
		WillMessage: []byte("offline"),
		WillRetain:  true,
	})
	if cliErr != nil {
		fmt.Println(cliErr)
		os.Exit(1)
	}

	// listen for commands
	subErr := cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("home/+/+/set"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					parts := strings.Split(string(topicName), "/")

					packet := noolite.Packet{}

					// detect device type
					switch parts[1] {
					case "noolite":
						packet.Mode = noolite.ModeTx
					case "nooolitef":
						packet.Mode = noolite.ModeFTx
					default:
						fmt.Println("expected noolite or noolitef device types but found:", parts[2])
						return
					}

					// get device channel
					channel, convErr := strconv.Atoi(parts[2])
					if convErr != nil || channel < 0 || channel > 63 {
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
							if err := nooDevice.Send(bind); err != nil {
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
						if err := nooDevice.Send(on); err != nil {
							fmt.Println(err)
							return
						}
					case "OFF": // turns device off
						off := packet
						off.Command = noolite.CmdOff
						if err := nooDevice.Send(off); err != nil {
							fmt.Println(err)
							return
						}
					default:
						// inform command does not supported
						fmt.Println("command does not supported:", command)
						return
					}

					fmt.Println("command", command, "sets on channel", channel)

					// send the current state
					if err := cli.Publish(&client.PublishOptions{
						QoS:       mqtt.QoS0,
						TopicName: []byte(fmt.Sprintf("home/%s/%s/state", parts[1], parts[2])),
						Message:   []byte(command),
					}); err != nil {
						fmt.Println(err)
					}
				},
			},
		},
	})
	if subErr != nil {
		fmt.Println(111, subErr)
		os.Exit(1)
	}
	fmt.Println("ready")

	// wait for the signal
	<-sigc
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}
