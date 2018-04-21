package main

import (
	"encoding/json"
	"fmt"
	"noolite/noolite"
	"os"
	"os/signal"
	"strings"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

// https://www.home-assistant.io/components/light.mqtt_json/
type Color struct {
	R uint8
	G uint8
	B uint8
	X float32
	Y float32
}
type Command struct {
	Brightness int
	ColorTemp  int `json:"color_temp"`
	Color      Color
	Effect     string
	State      string
	Transition uint8
	WhiteValue uint8 `json:"white_value"`
}

func main() {

	// handle interrupt signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// open noolite connected device
	device, nooErr := noolite.CreateDevice()
	if nooErr != nil {
		fmt.Println(nooErr)
		os.Exit(1)
	}
	defer device.Port.Close()

	// connect to MQTT server
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
			os.Exit(1)
		},
	})
	defer cli.Terminate()
	cliErr := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("example-client"),
	})
	if cliErr != nil {
		fmt.Println(cliErr)
		os.Exit(1)
	}

	// listen for commands
	subErr := cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("noolite/+/cmd"),
				QoS:         mqtt.QoS0,
				Handler: func(topicName, message []byte) {
					var command Command
					if err := json.Unmarshal(message, &command); err != nil {
						fmt.Println("error while decoding command:", err)
						return
					}
					parts := strings.Split(string(topicName), "/")
					channel := parts[1]
					fmt.Printf("%s, %+v", channel, command)
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
	device.Close()
	fmt.Println("disconnected")
	/*
		var mode, control, channel, command byte

		for _, input := range os.Args[1:] {
			if cmd, ok := modes[input]; ok == true {
				mode = cmd.value
				fmt.Println("adapter mode:", input)
				continue
			}
			if cmd, ok := txControls[input]; ok == true {
				fmt.Println("adapter state:", input)
				control = cmd.value
				continue
			}
			if cmd, ok := commands[input]; ok == true {
				fmt.Println("command:", input)
				command = cmd.value
				continue
			}
			possibleChannel, err := strconv.Atoi(input)
			if err == nil {
				fmt.Println("channel:", input)
				channel = byte(possibleChannel)
			}
		}

		if err := device.Send(mode, control, channel, command); err != nil {
			fmt.Println(error)
			os.Exit(1)

		}

		fmt.Println("command sent...")
	*/
}
