package main

import (
	"flag"
	"log"
	"noolite/noolite"
	"runtime"
	"strconv"
	"strings"

	"github.com/nats-io/go-nats"
)

func usage() {
	log.Fatalf("Usage: nats-req [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	// args := flag.Args()
	// if len(args) < 1 {
	// 	usage()
	// }

	device, createError := noolite.CreateDevice()
	if createError != nil {
		log.Fatalf("Can't init noolite: %v\n", createError)
	}

	nc, connectError := nats.Connect(*urls)
	if connectError != nil {
		log.Fatalf("Can't connect NATS: %v\n", connectError)
	}

	// emit driver events into queue
	go func() {
		for {
			input, err := device.Receive()
			if err != nil {
				log.Fatalf("Error while reading input buffer: %v\n", err)
			}
			// 2do: convert to nats-event and emit it
			log.Println("event:", input.Mode, input.Control, input.Channel, input.Command)
		}
	}()

	// noolite.send.{command}.{channel}.[address]
	// noolite.bind.on.{channel}.[address]
	// noolite.bind.off.{channel}.[address]
	// noolite.clear
	// noolite.clear.{channel}

	// noolite.legacy. ...

	//? noolite.state.{channel}

	// subsribe to incoming requests: "noolite.{command}.{channel}[.address]"
	nc.Subscribe("noolite.command.*.>", func(m *nats.Msg) {
		chunks := strings.SplitN(m.Subject, ".", 4)

		command := chunks[1]
		preset, exists := Commands[command]
		if exists != true {
			nc.Publish(m.Reply, []byte("+ERR: invalid command"))
			return
		}

		channel, convertErr := strconv.Atoi(chunks[2])
		if convertErr != nil {
			nc.Publish(m.Reply, []byte("+ERR: invalid channel"))
			return
		}
		if channel < 0 || channel > 63 {
			nc.Publish(m.Reply, []byte("+ERR: invalid channel (0..63)"))
		}

		device.Send(preset.Mode, preset.Control, byte(channel), preset.Command)
		// 2do: wait for success ?
		nc.Publish(m.Reply, []byte("+OK"))
	})

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Println("Service started")
	// fmt.Printf("%q\n", strings.SplitN("a,b,c", ",", 10))

	runtime.Goexit()
}
