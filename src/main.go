package main

import (

	"flag"
	"fmt"
	myhttp "github.com/shaojintian/load_balancer/src/http"
	"log"
	"net/http"

	"net/url"
	"strings"

)

// go run main.go --backends=http://localhost:3032,http://localhost:3033,http://localhost:3034
//测试： kbang http://localhost:3030

func main() {
	log.Println("start init--------------------------------")
	var serverPoll myhttp.ServerPool
	StartServer(serverPoll)
}



func StartServer(sp myhttp.ServerPool){
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



	//concurrency:heart beat check in other goroutine
	go sp.HeartBeatCheck()


	// init reverseProxy http handler
	server := http.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", 3030),
		Handler: http.HandlerFunc(sp.LoadBalance),
	}

	log.Println("start server.ListenAndServe--------------------------------")
	if err := server.ListenAndServe();err != nil {
		log.Fatal(err.Error())
	}


}

