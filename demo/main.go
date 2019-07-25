package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/factorysh/alertmanager-multitenant/multitenant"
)

func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url, _ := url.Parse("http://alertmanager:9093")
		proxy := httputil.NewSingleHostReverseProxy(url)

		r.URL.Host = url.Host
		r.URL.Scheme = url.Scheme
		r.Host = url.Host

		proxy.ServeHTTP(w, r)
		// client := &http.Client{}
		// r.URL.Host = url.Host
		// r.URL.Scheme = url.Scheme
		// r.RequestURI = ""
		// _, err := client.Do(r)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	})
	m := multitenant.Multitenant{
		JwtSecret: []byte("secret"),
	}
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", m.Multitenant(h)))
}
