package multitenant

import (
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

type Multitenant struct {
	JwtSecret []byte
}

type Data struct {
	Labels map[string]string `json:"labels"`
}

type Claims struct {
	Project string `json:"project"`
	jwt.StandardClaims
}

func (m *Multitenant) Multitenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data Data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		project := data.Labels["project"]
		// No project label
		if project == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// No jwt in header
		jwtStr := r.Header.Get("JWT")
		if jwtStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		// Parse the JWT string and store the result in `claims`.
		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
			return m.JwtSecret, nil
		})
		if err != nil {
			// Invalid signature
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// Bad JWT
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Invalid JWT
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
