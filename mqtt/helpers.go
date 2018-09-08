package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const brightnessStep = 0.419

func guessCommand(buf []byte) (command, payload string) {
	input := fmt.Sprintf("%s", buf)
	// check if brightness
	_, err := strconv.Atoi(input)
	if err == nil {
		return "BRIGHTNESS", input
	}
	rgbParts := strings.Split(input, ",")
	if len(rgbParts) == 3 { // looks like RGB
		return "RGB", input
	}
	return strings.ToUpper(input), ""

}

func validateByteRange(input string) (byte, error) {
	item, err := strconv.Atoi(input)
	if err != nil || item < 0 || item > 255 {
		return 0, errors.New("invalid range (0..255)")
	}
	return byte(item), nil
}

// converts 0..255 brightess scale into noolite specific
func scaleBrightness(baseScale byte) byte {
	// SN111-300: 43 .. 150 ?
	// return byte(brightnessStep*float64(baseScale) + 42)
	return byte(brightnessStep*float64(baseScale) + 42)
}
