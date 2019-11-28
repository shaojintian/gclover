package main

import (
	"flag"
	"fmt"
	myhttp"github.com/shaojintian/load_balancer/src/http"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	//"github.com/shaojintian/load_balancer/main"
)

func TestMain(m *testing.M) {
	log.Println("start init--------------------------------")
	var sp myhttp.ServerPool

	var serverList string
	//terminal read
	flag.StringVar(&serverList, "backends","","eg:http://localhost:3000,http://localhost:3001,...")
	flag.Parse()
	if len(serverList) == 0{
		log.Fatalln("invalid server list")
	}
	//do loadbalance
	urls := strings.Split(serverList, ",")
	for _, u := range urls{
		u,_ := url.Parse(u)
		sp.InitServerPool(u)
	}

	//show server pool
	for _, peer := range sp.Backends {
		fmt.Println(peer.URL.Host)
	}


	//concurrency:heart beat check in other goroutine
	//go sp.HeartBeatCheck()




	// init reverseProxy http handler
	server := http.Server{
		Addr: fmt.Sprintf("10.69.253.170:%d", 3030),
		Handler: http.HandlerFunc(sp.LoadBalance),
	}

	log.Println("start server.ListenAndServe--------------------------------")
	if err := server.ListenAndServe();err != nil {
		log.Fatal(err.Error())
	}



}
