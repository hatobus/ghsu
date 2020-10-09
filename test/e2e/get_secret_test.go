package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/hatobus/ghsu/test/fakeserver"
	"github.com/hatobus/ghsu/updator"
)

func prepareFakeServer(owner, repo string) *http.ServeMux {
	mux := http.NewServeMux()
	registerEndpoint := fmt.Sprintf("repos/%v/%v/actions/secrets", owner, repo)
	mux.HandleFunc(registerEndpoint, fakeserver.FakeGithubSecretHandler())
	return mux
}

func TestGetSecretFromServer(t *testing.T) {
	type testData struct {
		data       map[string]string
		wantStatus int
	}

	testCases := map[string]testData{
		"正常系": {
			data: map[string]string{
				"name": "hatobus",
				"age":  "22",
				"from": "japan",
			},
			wantStatus: http.StatusOK,
		},
	}

	client := updator.GithubClient{}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		m := prepareFakeServer(client.Owner, client.Repo)

		secretHandler := httptest.NewServer(m)
		u, err := url.Parse(secretHandler.URL)
		if err != nil {
			t.Fatal(err)
		}
		client.Client.BaseURL = u

		t.Cleanup(func() {
			secretHandler.Close()
		})
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			t.Log(tc.data)
		})
	}
}
