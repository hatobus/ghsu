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

	"github.com/google/go-github/github"
)

var mu sync.RWMutex
var secrets map[string]*github.Secret

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
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
	w.WriteHeader(http.StatusCreated)

	secret := new(github.EncryptedSecret)
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

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

	mu.Lock()
	secrets[secret.Name] = secretData
	mu.Unlock()
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, secretName := shiftPath(r.URL.Path)

	mu.RLock()
	secret, ok := secrets[secretName]
	mu.RUnlock()

	if !ok {
		http.Error(w, fmt.Sprintf("%v not found", secretName), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(secret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
