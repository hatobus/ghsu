package updator

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
)

var (
	ErrClientGithubKeyNotFound = xerrors.New("github token not found")
	ErrFailedToGetPublicKey    = xerrors.New("failed to get public key")
	ErrFailedToGetOwnerOrRepo  = xerrors.New("failed to get owner or repo")
)

type GithubClient struct {
	Token  string
	PublicKey *github.PublicKey
	Client *github.Client
}

func NewGithubClient() (*GithubClient, error) {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return nil, ErrClientGithubKeyNotFound
	}

	sts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})

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
		Token:  githubToken,
		Client: githubClient,
		PublicKey: pubKey,
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

		switch(key){
		case "user.name":
			owner = feature[1]
		case "remote.origin.url":
			repo = feature[1]
		default:
		}
	}

	return owner, repo, nil
}
