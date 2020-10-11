package fakeserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v31/github"
)

var mu sync.RWMutex
var secrets = map[string]*github.Secret{}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func getLastPathParameter(s string) string {
	p := strings.Split(s, "/")
	return p[len(p)-1]
}

func FakeGithubSecretHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getSecret(w, r)
			return
		case http.MethodPut:
			setSecret(w, r)
			return
		default:
			errStr := fmt.Sprintf("method: %v is not implemented", r.Method)
			http.Error(w, errStr, http.StatusNotImplemented)
			return
		}

	}
}

func setSecret(w http.ResponseWriter, r *http.Request) {
	var secret github.EncryptedSecret
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, secretPath := shiftPath(r.URL.Path)
	secret.Name = getLastPathParameter(secretPath)

	mu.Lock()
	s, ok := secrets[secret.Name]
	mu.Unlock()

	var ts github.Timestamp
	if ok {
		ts = s.CreatedAt
	} else {
		ts = github.Timestamp{Time: time.Now()}
	}

	secretData := &github.Secret{
		Name:      secret.Name,
		CreatedAt: ts,
		UpdatedAt: github.Timestamp{
			Time: time.Now(),
		},
	}

	log.Printf("secret: %v", secretData)

	mu.Lock()
	secrets[secret.Name] = secretData
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	_, secretPath := shiftPath(r.URL.Path)
	secretName := getLastPathParameter(secretPath)

	secret, ok := secrets[secretName]

	if !ok {
		http.Error(w, fmt.Sprintf("%v not found", secretName), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
