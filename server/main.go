package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr = flag.String("addr", ":6096", "Address of application")
	flag.Parse()

	r := newRoom()

	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
