//go:generate go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go
//go:generate env GOOS=js GOARCH=wasm go build -o ./build/app.wasm .
package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"strconv"

	"github.com/syumai/workers"
)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		msg := "Hello!"
		w.Write([]byte(msg))
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		io.Copy(w, bytes.NewReader(b))
	})

	http.HandleFunc("/image-placeholder", func(w http.ResponseWriter, r *http.Request) {
		width, err := strconv.ParseUint(r.URL.Query().Get("width"), 10, 16)
		if err != nil {
			http.Error(w, "invalid width passed", http.StatusBadRequest)
			return
		}
		height, err := strconv.ParseUint(r.URL.Query().Get("height"), 10, 16)
		if err != nil {
			http.Error(w, "invalid height passed", http.StatusBadRequest)
			return
		}

		if width > 2000 || height > 2000 {
			http.Error(w, "width or height too large.", http.StatusBadRequest)
			return
		}

		img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		gray := color.RGBA{R: 128, G: 128, B: 128, A: 255}
		for x := 0; x < int(width); x++ {
			for y := 0; y < int(height); y++ {
				img.Set(x, y, gray)
			}
		}

		w.Header().Set("Content-Type", "image/png")
		if err := png.Encode(w, img); err != nil {
			http.Error(w, "failed to encode image", http.StatusInternalServerError)
		}
	})

	workers.Serve(nil) // use http.DefaultServeMux
}
