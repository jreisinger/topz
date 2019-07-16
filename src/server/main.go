package main

import (
	"log"
	"net/http"

	"topz"
)

func main() {
	http.HandleFunc("/topz", topz.HandleTopz)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
