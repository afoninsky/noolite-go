package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/afoninsky/noolite-go/noolite"
	"github.com/spf13/viper"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {

	// configuration
	viper.SetDefault("mqtt.host", "localhost:1883")
	viper.BindEnv("mqtt.host", "MQTT_HOST")

	// handle interrupt signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// open noolite connected device
	nooDevice, nooErr := noolite.CreateDevice()
	if nooErr != nil {
		log.Fatalln(nooErr)
	}
	defer nooDevice.Close()

	// connect to MQTT server
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Fatalln(err)
		},
	})
	defer cli.Terminate()

	if err := cli.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      viper.GetString("mqtt.host"),
		ClientID:     []byte(clientID),
		CleanSession: true,
		WillTopic:    []byte(willTopic),
		WillMessage:  []byte(willOfflineMessage),
		WillRetain:   true,
	}); err != nil {
		log.Fatalln(err)
	}

	server := &Server{noolite: &nooDevice, mqtt: cli}

	// listen for incoming commands
	if err := cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(fmt.Sprintf(setTopicPattern, "+", "+")),
				QoS:         mqtt.QoS0,
				Handler:     server.messageHandler,
			},
		},
	}); err != nil {
		log.Fatalln(err)
	}

	// wait for events in bus
	// go func() {
	// 	for {
	// 		input, rcvErr := server.noolite.Receive()
	// 		if rcvErr != nil {
	// 			log.Println(rcvErr)
	// 			return
	// 		}
	// 		// input.Command - command
	// 		// input.Channel - channel
	// 		// input.Data: [1,0,0,0] on fail (device is diabled)
	// 		// input.Data: [1,2,0,0] on fail (device is enabled but command not complete)
	// 		fmt.Println("command:", input.Command)
	// 		fmt.Println("channel:", input.Channel)
	// 		fmt.Println("data:", input.Data)
	// 		fmt.Println("datafmt:", input.DataFormat)
	// 	}
	// }()

	// send readiness message into status topic
	if err := cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		Retain:    true,
		TopicName: []byte(willTopic),
		Message:   []byte(willOnlineMessage),
	}); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Ready to accept incoming connections")

	<-sigc
	// send offline message
	if err := cli.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		Retain:    true,
		TopicName: []byte(willTopic),
		Message:   []byte(willOfflineMessage),
	}); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second)

	// disconnect
	if err := cli.Disconnect(); err != nil {
		log.Fatalln(err)
	}

}
