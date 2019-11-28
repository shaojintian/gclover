package http

import (
	"github.com/shaojintian/load_balancer/src/core"
	"log"
	"net"
	"time"
)

//heart beat to refresh backend

func (sp *ServerPool)HeartBeatCheck(){
	// 2 mins to refresh all backend status
	var count int
	count = 0
	for{

		select{
		case <-time.After(time.Second *20) :
			log.Printf("[HeartBeatCheck]:----- %d round----- start",count)
			for _,peer := range sp.Backends {
				alive := backendHeartBeatAlive(peer,count)
				peer.SetAlive(alive)
			}
			log.Printf("[HeartBeatCheck]:----- %d round----- completed",count)
		}
		count++
	}

}

func backendHeartBeatAlive(peer *core.Backend,count int) bool {
	timeout := 2 * time.Second
	//url.Host == ip:port
	conn, err := net.DialTimeout("tcp",peer.URL.Host,timeout)
	if err != nil {
		log.Printf("[HeartBeatCheck]:----- %d round----- %s [dead]\n",count,peer.URL.Host)
		return false
	}
	// has connection
	connErr := conn.Close()
	if connErr != nil{
		panic(connErr.Error())
	}
	log.Printf("[HeartBeatCheck]:----- %d round----- %s [alive]\n",count,peer.URL.Host)
	return true

}



