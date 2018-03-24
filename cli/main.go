package main

import (
	"fmt"
	"noolite/noolite"
	"os"
	"strconv"
)

func main() {

	device, error := noolite.CreateDevice()
	if error != nil {
		fmt.Println(error)
		os.Exit(1)
	}
	defer device.Port.Close()

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

}
