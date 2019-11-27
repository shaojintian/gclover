package core

import (
	myHttp "github.com/shaojintian/load_balancer/src/http"
	"net/http"
)

func (sp *myHttp.ServerPool)LoadBalance(rw http.ResponseWriter,req *http.Request){
	perr := sp.GetNextPeer()
	
}
