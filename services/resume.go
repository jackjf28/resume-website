package services

import (
	"context"
	"fmt"
	"time"
	"github.com/jackjf28/resume-website/interfaces"
	"github.com/jackjf28/resume-website/utils"
)

type ResumeService struct {
	githubClient interfaces.GitHubClient
	cache        *utils.TTLCache[string, []byte]
}

func NewResumeService(githubClient interfaces.GitHubClient) *ResumeService {
	return &ResumeService{
		githubClient: githubClient,
		cache:        utils.NewTTLCache[string, []byte](),
	}
}

func (rs *ResumeService) GetResumePDF(ctx context.Context) ([]byte, error) {
	if pdfData, ok := rs.cache.Get("resume"); ok {
		return pdfData, nil
	}
	pdfData, err := rs.githubClient.GetPDFBytes(ctx, "jackjf28", "resume", "jack_farrell_resume.pdf")
	if err != nil {
		return nil, fmt.Errorf("fetching resume from GitHub: %v", err)
	}
	rs.cache.Set("resume", pdfData, time.Minute * 10)
	return pdfData, nil
}
