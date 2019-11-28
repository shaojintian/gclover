package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
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
func (sp *ServerPool)DoLoadBalance(serverUrl *url.URL) {
	reverseProxy := httputil.NewSingleHostReverseProxy(serverUrl)

	//add backend in sp
	sp.AddBackend(serverUrl, true,reverseProxy)
	log.Printf("Cofigured this backend:%s\n", serverUrl)

	//  to handle ReverseProxy error in callback func ReverseProxy.ErrorHandler()
	//  async function ,only run this when callback and don't do bellow code and other all operations
	reverseProxy.ErrorHandler = func(rw http.ResponseWriter,req *http.Request,err error){
		log.Printf("[host: %s],%s",serverUrl.Host,err.Error())
		// handle retry for this peer
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

		//retry time >= 3 to kill this URL
		//and LoadBalance() to select next peer
		sp.MarkPeerStatus(serverUrl,false)
		//  记录某一个请求次数，handle各种情况
		attemps := GetAttemptsFromCtx(req)
		ctx := context.WithValue(req.Context(),Attempts,attemps+1)
		sp.LoadBalance(rw,req.WithContext(ctx))

	}
}