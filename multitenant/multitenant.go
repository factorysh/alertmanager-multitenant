package multitenant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		var datas []*Data
		// Read the Body for clonning it
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Reset the Body as Before
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		// Unmarshal
		err = json.NewDecoder((bytes.NewReader(body))).Decode(&datas)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		project := datas[0].Labels["project"]
		// No project label
		if project == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// No jwt in header
		jwtStr := r.Header.Get("JWT")
		if jwtStr == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// Good path
		if r.URL.Path != "/api/v2/alerts" || r.URL.RawQuery != "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		claims := &Claims{}
		// Parse the JWT string and store the result in `claims`.
		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return m.JwtSecret, nil
		})
		if err != nil {
			// Bad JWT
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// Invalid JWT
		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
