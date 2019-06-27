package controller

import (
	"log"
	"net/http"
	"strconv"
)

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeBytes(w http.ResponseWriter, b []byte, format string) {
	w.Header().Set("Content-Type", "image/"+format)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	if _, err := w.Write(b); err != nil {
		log.Println("unable to write image.")
	}
}
