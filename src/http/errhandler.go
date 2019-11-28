package http

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (sp *ServerPool)reverseProxyErrHandler(reverseProxy *httputil.ReverseProxy,serverUrl *url.URL){
	//  to handle ReverseProxy error in callback func ReverseProxy.ErrorHandler()
	//  async function ,only run this when callback and don't do bellow code and other all operations
	reverseProxy.ErrorHandler = func(rw http.ResponseWriter,req *http.Request,err error){
		log.Printf("reverseProxyErrHandler: [host: %s],error--------:%s",serverUrl.Host,err.Error())
		// handle retry for this peer
		retries := GetRetryFromCtx(req)
		log.Printf("Req.Url %s retry %d times",req.URL.Host,retries)
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
		log.Printf("Req.Url: %s attempt %d times",req.URL.Host,attemps)
		ctx := context.WithValue(req.Context(),Attempts,attemps+1)
		sp.LoadBalance(rw,req.WithContext(ctx))

	}
}
