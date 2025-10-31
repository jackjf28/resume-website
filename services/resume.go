package services

import (
	"context"
	"fmt"

	"github.com/jackjf28/resume-website/interfaces"
)

type ResumeService struct {
	githubClient interfaces.GitHubClient
	cache        map[string][]byte
}

func NewResumeService(githubClient interfaces.GitHubClient) *ResumeService {
	return &ResumeService{
		githubClient: githubClient,
		cache:        make(map[string][]byte),
	}
}

func (rs *ResumeService) GetResumePDF(ctx context.Context) ([]byte, error) {
	if pdfData, ok := rs.cache["resume"]; ok {
		return pdfData, nil
	}
	pdfData, err := rs.githubClient.GetPDFBytes(ctx, "jackjf28", "resume", "jack_farrell_resume.pdf")
	if err != nil {
		return nil, fmt.Errorf("fetching resume from GitHub: %v", err)
	}
	rs.cache["resume"] = pdfData
	return pdfData, nil
}
