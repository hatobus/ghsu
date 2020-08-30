package updator

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/xerrors"
	"os"
)

var (
	ErrClientGithubKeyNotFound = xerrors.New("github token not found")
	ErrFailedToGetPublicKey = xerrors.New("failed to get public key")
)

type GithubClient struct {
	Token string
	PubKey string
	Owner string
	Repo string
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

	return &GithubClient{
		Token: githubToken,
		Client: githubClient,
	}, nil
}

func (g *GithubClient) GetPublicKeyFromGitHub(owner, repo string) (string, error) {
	panic("not impl")
}
