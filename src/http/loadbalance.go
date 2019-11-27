package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

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

func (sp *ServerPool)DoLoadBalance(serverUrl *url.URL) {
	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)
	//  to handle ReverseProxy error in callback func ReverseProxy.ErrorHandler()
	//  async function ,only run this when callback and don't do bellow code and other all operations
	reverseProxy.ErrorHandler = func(rw http.ResponseWriter,req *http.Request,err error){
		log.Printf("[host: %s],%s",serverUrl.Host,err.Error())
		// handle retry
		retries := GetRetryFromCtx(req)
		if retries < 3 {
			// timeout 10ms to retry
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(req.Context(),Retry,retries+1)
				reverseProxy.ServeHTTP(rw,req.WithContext(ctx))
				return
			}
		}

		//retry time >= 3 to kill this peer
		sp.MarkPeerStatus(serverUrl,false)


	}
}