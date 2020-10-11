package updator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

var (
	ErrClientGithubKeyNotFound = xerrors.New("github token not found")
	ErrFailedToGetPublicKey    = xerrors.New("failed to get public key")
	ErrFailedToGetOwnerOrRepo  = xerrors.New("failed to get owner or repo")
)

type GithubClient struct {
	Token     string
	PublicKey *github.PublicKey
	Client    *github.Client
	Owner     string
	Repo      string
}

func NewGithubClient() (*GithubClient, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, ErrClientGithubKeyNotFound
	}

	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)

	oauthClient := oauth2.NewClient(context.Background(), sts)
	githubClient := github.NewClient(oauthClient)

	owner, repo, err := getOwnersAndRepoFromCurrentGitFile()
	if err != nil {
		return nil, err
	} else if owner == "" || repo == "" {
		return nil, ErrFailedToGetOwnerOrRepo
	}

	pubKey, _, err := githubClient.Actions.GetPublicKey(context.Background(), owner, repo)
	if err != nil {
		return nil, err
	}

	return &GithubClient{
		Token:     githubToken,
		Client:    githubClient,
		PublicKey: pubKey,
		Owner:     owner,
		Repo:      repo,
	}, nil
}

func getOwnersAndRepoFromCurrentGitFile() (string, string, error) {
	out, err := exec.Command("git", "config", "--list").Output()
	if err != nil {
		return "", "", err
	}

	configs := strings.Split(string(out), "\n")

	var owner, repo string

	for _, s := range configs {
		feature := strings.Split(s, "=")
		key := feature[0]

		switch key {
		case "user.name":
			owner = feature[1]
		case "remote.origin.url":
			c := strings.TrimRight(feature[1], ".git")
			s := strings.Split(c, "/")
			repo = s[len(s)-2]
		default:
		}
	}

	return owner, repo, nil
}

func (gc *GithubClient) GenerateEncryptedSecret(data map[string]string) ([]*github.EncryptedSecret, error) {
	pk, err := gc.getRawPublicKey()
	if err != nil {
		return nil, err
	}

	secrets := make([]*github.EncryptedSecret, 0, len(data))

	for name, value := range data {
		secrets = append(secrets, &github.EncryptedSecret{
			Name:           name,
			KeyID:          *gc.PublicKey.KeyID,
			EncryptedValue: encryptSodium(value, pk),
		})
	}

	return secrets, nil
}

func (gc *GithubClient) ShowUpSetSecrets() error {
	ctx := context.Background()
	secrets, res, err := gc.Client.Actions.ListSecrets(ctx, gc.Owner, gc.Repo, nil)
	if err != nil {
		log.Println(res)
		return err
	}

	fmt.Printf("github secret %v already set\n", secrets.TotalCount)

	for _, s := range secrets.Secrets {
		fmt.Printf("variable: %v\n", s.Name)
		fmt.Printf("created at: %v\n", s.CreatedAt)
		fmt.Printf("last update: %v\n", s.UpdatedAt)
	}

	return nil
}

func (gc *GithubClient) ExistRepoSecret(owner, repo, name string) bool {
	ctx := context.Background()
	_, res, err := gc.Client.Actions.GetSecret(ctx, owner, repo, name)
	if res != nil && res.StatusCode != http.StatusOK {
		return false
	} else if err != nil {
		log.Println(err)
		return false
	}
	return true
}
