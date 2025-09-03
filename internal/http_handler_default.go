package internal

import (
	"net/http"
)

func DefaultHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	_, err := resp.Write([]byte("This is SheepsTor"))
	if err != nil {
		Log.Error(err.Error())
		return
	}
}
