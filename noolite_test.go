package noolite

import (
	"testing"
)

func TestDebug(t *testing.T) {
	_, createError := CreateDevice()
	if createError != nil {
		t.Error(createError)
	}

}
