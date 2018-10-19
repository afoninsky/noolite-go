package main

import (
	"fmt"
	"strconv"
	"strings"
)

// SN111-300 has brightness from 43.. 150
// this const helps convert standart 0..255 scale to noolite specific
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

func validateByteRange(input string, min, max int) (byte, error) {
	item, err := strconv.Atoi(input)
	if err != nil || item < min || item > max {
		return 0, fmt.Errorf("invalid range (%v..%v)", min, max)
	}
	return byte(item), nil
}

func scaleBrightness(baseScale byte) byte {
	return byte(brightnessStep*float64(baseScale) + 42)
}
