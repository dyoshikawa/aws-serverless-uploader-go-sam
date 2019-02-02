package main

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"testing"
)

func TestResizePng(t *testing.T) {
	t.Run("Resize png", func(t *testing.T) {
		file, err := os.Open("450-150.png")
		if err != nil {
			t.Fatal("ERROR: " + err.Error())
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			t.Fatal("ERROR: " + err.Error())
		}
		resizedBytes, err := resizePng(fileBytes)
		if err != nil {
			t.Fatal("ERROR: " + err.Error())
		}

		r := bytes.NewReader(resizedBytes)
		im, _, err := image.DecodeConfig(r)
		if err != nil {
			t.Fatal("ERROR: " + err.Error())
		}

		if im.Width != 300 {
			t.Fatal("FAILED: width of resized is not 300.")
		}
	})
}
