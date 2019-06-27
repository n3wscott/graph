package controller

import (
	"bytes"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/codec"
	"image"
	"image/jpeg"
	"log"
	"net/http"
)

func (c *Controller) Favicon(w http.ResponseWriter, r *http.Request) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	writeBytes(w, buffer.Bytes(), "jpg")
}
