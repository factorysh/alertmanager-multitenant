package multitenant

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func yes(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestMiddleware(t *testing.T) {
	m := &Multitenant{}
	ts := httptest.NewServer(m.Multitenant(http.HandlerFunc(yes)))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	
	fmt.Println(res, err)
}
