package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	fmt.Println("hello")
}

func init(){
	var port int
	//init ReverseProxy
	targetUrl,_ := url.Parse("http://localhost:8080")
	//先快速开发，再处理err  最后在补充处理error
	reverseProxy := httputil.NewSingleHostReverseProxy(targetUrl)

	// init reverseProxy http handler
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(reverseProxy.ServeHTTP),
	}

}