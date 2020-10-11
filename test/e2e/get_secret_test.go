package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v31/github"

	"github.com/hatobus/ghsu/test/fakeserver"
	"github.com/hatobus/ghsu/updator"
)

func prepareFakeServer(owner, repo string) *http.ServeMux {
	mux := http.NewServeMux()
	registerEndpoint := fmt.Sprintf("/api/v3/repos/%v/%v/actions/secrets/", owner, repo)
	mux.HandleFunc(registerEndpoint, fakeserver.FakeGithubSecretHandler())
	return mux
}

func toPtr(s string) *string {
	return &s
}

func TestGetSecretFromServer(t *testing.T) {
	type testData struct {
		data    map[string]string
		isExist map[string]bool
	}

	defaultData := map[string]string{
		"name": "hatobus",
		"age":  "22",
		"from": "japan",
	}

	testCases := map[string]testData{
		"正常系": {
			data: defaultData,
			isExist: map[string]bool{
				"name": true,
				"age":  true,
				"from": true,
			},
		},
		"存在しないものも入っている": {
			data: defaultData,
			isExist: map[string]bool{
				"name":   true,
				"hobby":  false,
				"weight": false,
			},
		},
	}

	var client *updator.GithubClient

	if os.Getenv("GITHUB_ACTIONS") != "true" {
		owner := "hatobus"
		repo := "ghsu"

		m := prepareFakeServer(owner, repo)

		secretHandler := httptest.NewServer(m)

		httpClient := http.DefaultClient

		githubClient, err := github.NewEnterpriseClient(secretHandler.URL, secretHandler.URL, httpClient)
		if err != nil {
			t.Fatal(err)
		}

		client = &updator.GithubClient{
			Owner: owner,
			Repo:  repo,
			PublicKey: &github.PublicKey{
				Key:   toPtr(strings.Repeat("a", 32)),
				KeyID: toPtr("key-id"),
			},
			Client: githubClient,
		}

		t.Cleanup(func() {
			secretHandler.Close()
		})
	} else {
		var err error
		client, err = updator.NewGithubClient()
		if err != nil {
			t.Fatal(err)
		}
	}

	for testName, tc := range testCases {
		t.Run(testName, func(t *testing.T) {
			encrypted, err := client.GenerateEncryptedSecret(tc.data)
			if err != nil {
				t.Fatal(err)
			}

			// add secret data mock server
			for _, secret := range encrypted {
				ctx := context.Background()
				_, err := client.Client.Actions.CreateOrUpdateSecret(ctx, client.Owner, client.Repo, secret)
				if err != nil {
					t.Log(err)
					t.Fatalf("Key: %v, registration error", secret.Name)
				}
			}

			for name, exist := range tc.isExist {
				res := client.ExistRepoSecret(client.Owner, client.Repo, name)

				if diff := cmp.Diff(exist, res); diff != "" {
					t.Fatalf("invalid ExistRepoSecret response, diff: %v", diff)
				}
			}
		})
	}
}
