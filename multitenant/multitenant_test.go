package multitenant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test_struct struct {
	Body string
}

func yes(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var t test_struct
	json.Unmarshal(body, &t)
	log.Println(t.Body)
	// mockup alertmanager POST alert
	w.WriteHeader(http.StatusOK)
}

func TestMiddleware(t *testing.T) {
	m := &Multitenant{
		JwtSecret: []byte("secret"),
	}
	ts := httptest.NewServer(m.Multitenant(http.HandlerFunc(yes)))
	defer ts.Close()
	// Simple GET request
	res, err := http.Get(ts.URL)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "the status code should be StatusBadRequest")
	assert.NoError(t, err)
	// POST request without project label
	res, err = http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(`{"labels": {"blabla": "test"}}`)))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "the status code should be StatusBadRequest")
	assert.NoError(t, err)
	// POST request with project label but without JWT header
	res, err = http.Post(ts.URL, "application/json", bytes.NewBuffer([]byte(`{"labels": {"project": "test"}}`)))
	assert.Equal(t, http.StatusForbidden, res.StatusCode, "the status code should be StatusForbidden")
	assert.NoError(t, err)
	// Bad api path
	client := &http.Client{}
	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer([]byte(`{"labels": {"project": "test"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o")
	res, err = client.Do(req)
	assert.Equal(t, http.StatusForbidden, res.StatusCode, "the status code should be StatusForbidden")
	assert.NoError(t, err)
	// Path with param in url
	client = &http.Client{}
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v2/alerts?active=true", ts.URL), bytes.NewBuffer([]byte(`{"labels": {"project": "test"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o")
	res, err = client.Do(req)
	assert.Equal(t, http.StatusForbidden, res.StatusCode, "the status code should be StatusForbidden")
	assert.NoError(t, err)
	// Bad JWT
	client = &http.Client{}
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v2/alerts", ts.URL), bytes.NewBuffer([]byte(`{"labels": {"project": "test"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("JWT", "very_bad_jwt")
	res, err = client.Do(req)
	assert.Equal(t, http.StatusForbidden, res.StatusCode, "the status code should be StatusForbidden")
	assert.NoError(t, err)
	// Bad Method (RS256 instead of HS*)
	client = &http.Client{}
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v2/alerts", ts.URL), bytes.NewBuffer([]byte(`{"labels": {"project": "test"}}`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("JWT", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTMzNywidXNlcm5hbWUiOiJqb2huLmRvZSJ9.5czc5hsfIUkPpqKGVKOJlIKiSwyh1pezJESJU5DVCa0")
	res, err = client.Do(req)
	assert.Equal(t, http.StatusForbidden, res.StatusCode, "the status code should be StatusForbidden")
	assert.NoError(t, err)
	// Normal Case
	client = &http.Client{}
	jsonStr := []byte(`{"labels": {"project": "test"}}`)
	req, err = http.NewRequest("POST", fmt.Sprintf("%s/api/v2/alerts", ts.URL), bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("JWT", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o")
	res, err = client.Do(req)
	assert.Equal(t, http.StatusOK, res.StatusCode, "the status code should be StatusOK")
	assert.NoError(t, err)
	fmt.Println("Second :", res, err)
}
