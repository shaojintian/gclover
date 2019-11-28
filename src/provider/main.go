package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	// init reverseProxy http handler
	provider := http.Server{
		Addr: fmt.Sprintf("10.69.253.170:%d", 4001),
		//Handler: http.HandlerFunc(sp.LoadBalance),
	}

	log.Println("start provider[%s].ListenAndServe--------------------------------",provider.Addr)
	if err := provider.ListenAndServe();err != nil {
		log.Fatal(err.Error())
	}
}
