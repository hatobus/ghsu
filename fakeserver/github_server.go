package fakeserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/google/go-github/github"
)

var mu sync.RWMutex
var secrets map[string]string

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func FakeGithubRegisterSecretHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := new(github.EncryptedSecret)
		if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		mu.Lock()
		secrets[secret.Name] = secret.EncryptedValue
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
		return
	}
}

func FakeGithubSecretResponseHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, secretName := shiftPath(r.URL.Path)

		mu.RLock()
		_, ok := secrets[secretName]
		mu.RUnlock()

		if !ok {
			http.Error(w, fmt.Sprintf("%v not found", secretName), http.StatusBadRequest)
			return
		}

		s := &github.Secret{
			Name: secretName,
		}

		b, err := json.Marshal(s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}
