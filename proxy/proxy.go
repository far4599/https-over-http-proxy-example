package proxy

import (
	"io"
	"net/http"
)

func HTTPProxyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "CONNECT" {
		http.Error(w, "CONNECT not implemented", http.StatusInternalServerError)
		return
	}

	r.RequestURI = ""
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for header := range w.Header() {
		resp.Header.Add(header, w.Header().Get(header))
	}

	io.Copy(w, resp.Body)
}
