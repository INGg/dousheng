package util

import "testing"

func TestVideo2Image(t *testing.T) {
	fileName := "doudoudou"
	err := SaveImageFromVideo(fileName, true)
	if err != nil {
		return
	}
}
