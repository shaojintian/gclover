package http

import "net/http"

func GetRetryFromCtx(req *http.Request)  int{
	if retry,ok := req.Context().Value();ok{
		return retry
	}

	

}
