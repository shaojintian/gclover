package http

import "net/http"


func GetRetryFromCtx(req *http.Request)  int{
	if retry,ok := req.Context().Value(Retry).(int);ok{
		return retry
	}

	return 0

}

func GetAttemptsFromCtx(req *http.Request) int{
	if attempts,ok := req.Context().Value(Attempts).(int);ok{
		return attempts
	}
	return 0
}
