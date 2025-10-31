package github

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jackjf28/resume-website/utils"
)

const base64Encoding string = "base64"

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

type FileContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        uint   `json:"size"`
	URL         string `json:"url"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.github.com",
		token:      token,
	}
}

func (c *Client) GetFileContent(ctx context.Context, owner, repo, path string) (*FileContent, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", c.baseURL, owner, repo, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	c.setDefaultRequestHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	fileContent, err := utils.Decode[FileContent](resp)
	if err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}
	return &fileContent, nil
}

func (c *Client) GetPDFBytes(ctx context.Context, owner, repo, path string) ([]byte, error) {
	fileContent, err := c.GetFileContent(ctx, owner, repo, path)
	if err != nil {
		return nil, err
	}

	if fileContent.Encoding != base64Encoding {
		return nil, fmt.Errorf("expected base64 encoding, got %s", fileContent.Encoding)
	}

	pdfData, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return nil, fmt.Errorf("error decoding pdf data: %v", err)
	}

	return pdfData, nil
}

func (c *Client) setDefaultRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/vnd.github.object")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
}
