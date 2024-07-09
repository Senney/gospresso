package main

import (
	"net/http"

	"github.com/Senney/gospresso"
)

func main() {
	r := gospresso.NewMux()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	r.Get("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bar!"))
	})

	http.ListenAndServe(":3333", r)
}
