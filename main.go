package main

import (
	"log"
	"net/http"

	"github.com/far4599/https-over-http-proxy-example/proxy"
)

func main() {
	simpleProxyHandler := http.HandlerFunc(proxy.HTTPProxyHandlerFunc)
	if err := http.ListenAndServe(":44444", simpleProxyHandler); err != nil {
		log.Fatalf("server failed to listen to port 0.0.0.0:44444: '%v'", err)
	}
}
