package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

var jwtKey = []byte("secret")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	// get the token
	http.HandleFunc("/getjwt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		var creds Credentials
		// Get the JSON body and decode into credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		expectedPassword, ok := users[creds.Username]
		if !ok || expectedPassword != creds.Password {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		expirationTime := time.Now().Add(time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, tokenString)
	})
	//
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		jwtStr := r.Header.Get("JWT")
		if jwtStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

	})
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
