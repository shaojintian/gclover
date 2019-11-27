package http

import "net/http"

func (sp *ServerPool) LoadBalance(rw http.ResponseWriter,req *http.Request){
	peer := sp.GetNextPeer()
	if peer == nil {
		http.Error(rw,"no alive backend peer",http.StatusServiceUnavailable)
	}else{
		peer.ReverseProxy.ServeHTTP(rw,req)
		// to handle err in callback func ReverseProxy.ErrorHandler()
		return
	}
}
