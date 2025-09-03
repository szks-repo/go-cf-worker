//go:generate go run github.com/syumai/workers/cmd/workers-assets-gen -mode=go
//go:generate env GOOS=js GOARCH=wasm go build -o ./build/app.wasm .
package main

import (
	"bytes"
	"io"
	"net/http"

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

	workers.Serve(nil) // use http.DefaultServeMux
}
