package fakeserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"

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

func verifyAuthToken(r *http.Request) (int, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if diff := cmp.Diff(string(b), os.Getenv("GITHUB_TOKEN")); diff != "" {
		return http.StatusBadRequest, fmt.Errorf("Authentication token not valid, diff: %v\n", diff)
	}

	return http.StatusOK, nil
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
	code, err := verifyAuthToken(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), code)
	}

	w.WriteHeader(code)

	secret := new(github.EncryptedSecret)
	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	mu.Lock()
	secrets[secret.Name] = secret.EncryptedValue
	mu.Unlock()

	return
}

func getSecret(w http.ResponseWriter, r *http.Request) {
	code, err := verifyAuthToken(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), code)
	}

	w.WriteHeader(code)

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

	w.Write(b)
	return
}
