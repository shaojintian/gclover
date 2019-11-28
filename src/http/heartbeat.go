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
	for{
		select{
		case <-time.After(time.Minute*2) :
			for _,peer := range sp.backends{
				alive := backendHeartBeatAlive(peer)
				peer.SetAlive(alive)
			}

		}
	}

}

func backendHeartBeatAlive(peer *core.Backend) bool {
	timeout := 2 * time.Second
	//url.Host == ip:port
	conn, err := net.DialTimeout("tcp",peer.URL.Host,timeout)
	if err != nil {
		log.Printf("%s cannot reach\n",peer.URL.Host)
		return false
	}
	// has connection
	connErr := conn.Close()
	if connErr != nil{
		panic(connErr.Error())
	}

	return true

}



