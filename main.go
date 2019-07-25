package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/factorysh/alertmanager-multitenant/multitenant"
)

func main() {
	am_url := os.Getenv("AM_ADDRESS")
	if am_url == "" {
		log.Println("AM_URL env variable (alertmanager url) is not set. Default: http://0.0.0.0:9093")
		am_url = "http://0.0.0.0:9093"
	}
	listen_addr := os.Getenv("LISTEN_ADDRESS")
	if listen_addr == "" {
		log.Println("LISTEN_ADDRESS env variable is not set. Default: 0.0.0.0:9000")
		listen_addr = "0.0.0.0:9000"
	}
	secret := os.Getenv("SIGNATURE")
	if secret == "" {
		log.Println("SIGNATURE env variable (alertmanager address) is not set. Default: \"secret\"")
		secret = "secret"
	}
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, _ := url.Parse(am_url)
		proxy := httputil.NewSingleHostReverseProxy(url)
		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Host = url.Host

		proxy.ServeHTTP(w, r)
	})
	m := multitenant.Multitenant{
		JwtSecret: []byte(secret),
	}
	log.Fatal(http.ListenAndServe(listen_addr, m.Multitenant(h)))
}
