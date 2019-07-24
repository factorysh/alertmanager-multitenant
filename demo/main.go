package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/factorysh/alertmanager-multitenant/multitenant"
)

func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host)
		client := &http.Client{}
		r.URL.Host = "192.168.99.105:9093"
		_, err := client.Do(r)
		if err != nil {
			log.Fatal(err)
		}
	})
	m := multitenant.Multitenant{
		JwtSecret: []byte("secret"),
	}
	log.Fatal(http.ListenAndServe("127.0.0.1:9000", m.Multitenant(h)))
}
