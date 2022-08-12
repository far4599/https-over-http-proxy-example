package main

import (
	"net/http"

	"github.com/far4599/https-over-http-proxy-example/proxy"
)

func main() {
	simpleProxyHandler := http.HandlerFunc(proxy.HTTPProxyHandlerFunc)
	http.ListenAndServe(":44444", simpleProxyHandler)
}
