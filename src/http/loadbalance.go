package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//handle  one http attempt
func (sp *ServerPool) LoadBalance(rw http.ResponseWriter,req *http.Request){
	if GetAttemptsFromCtx(req) > len(sp.backends){
		log.Printf("Max attemps attached: %s",req.URL)
		http.Error(rw,"this req to much attemps",http.StatusServiceUnavailable)
		return
	}
	peer := sp.GetNextPeer()
	if peer == nil {
		http.Error(rw,"no alive backend peer",http.StatusServiceUnavailable)
	}else{
		peer.ReverseProxy.ServeHTTP(rw,req)
		// to handle err in callback func ReverseProxy.ErrorHandler()
		return
	}
}

//handle one specific url
// init backends
func (sp *ServerPool)InitServerPool(serverUrl *url.URL) {
	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)

	//add backend in sp
	sp.AddBackend(serverUrl, true,reverseProxy)
	log.Printf("Cofigured this backend:%s\n", serverUrl)

	//callback
	sp.reverseProxyErrHandler(reverseProxy, serverUrl)

}