package interfaces

import (
	"context"

	"github.com/jackjf28/resume-website/github"
)

type GitHubClient interface {
	GetPDFBytes(ctx context.Context, owner, repo, path string) ([]byte, error)
}

var _ GitHubClient = (*github.Client)(nil)
