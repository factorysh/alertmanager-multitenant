package multitenant

import "net/http"

type Multitenant struct {
}

func (m *Multitenant) Multitenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		next.ServeHTTP(w, r)
	})
}
