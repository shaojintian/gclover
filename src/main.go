package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//先快速开发，再处理err  最后在补充处理error
func main() {
	fmt.Println("hello")
}

func init(){
	var port int
	//init ReverseProxy
	targetUrl,_ := url.Parse("http://localhost:8080")

	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)
		//  to handle ReverseProxy error in callback func ReverseProxy.ErrorHandler()
		//  async function ,only run this when callback and don't do bellow code and other all operations
		reverseProxy.ErrorHandler = func(rw http.ResponseWriter,req *http.Request,err error){
			log.Printf("[host: %s],%s",targetUrl.Host,err.Error())


		}


	// init reverseProxy http handler
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(reverseProxy.ServeHTTP),
	}

}