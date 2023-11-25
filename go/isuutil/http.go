package isuutil

import (
	"net/http"
	"time"
)

func InitializeDefaultHTTPClient() {
	// DefaultTransportの制限を変更する場合。
	// DefaultTransportはhttp.DefaultClientのTransportとして使われる。
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 200
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 500
	http.DefaultTransport.(*http.Transport).IdleConnTimeout = 120 * time.Second
}
