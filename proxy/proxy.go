package proxy

import (
	"io"
	"net"
	"net/http"
	"time"
)

func HTTPProxyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		serveCONNECT(w, r)
		return
	}

	r.RequestURI = ""
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	for header := range resp.Header {
		w.Header().Set(header, resp.Header.Get(header))
	}

	io.Copy(w, resp.Body)
}

func serveCONNECT(w http.ResponseWriter, r *http.Request) {
	targetConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	go transmit(targetConn, clientConn)
	go transmit(clientConn, targetConn)
}

func transmit(from io.ReadCloser, to io.WriteCloser) {
	defer from.Close()
	defer to.Close()

	io.Copy(to, from)
}
