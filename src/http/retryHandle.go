package http

import "net/http"


func GetRetryFromCtx(req *http.Request)  int{
	if retry,ok := req.Context().Value(Retry).(int);ok{
		return retry
	}

	return 0

}
