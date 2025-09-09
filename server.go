package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const jsonContentType string = "application/json"

type Server struct {
	http.Handler
	ctx context.Context
}

func NewServer(context context.Context) http.Handler {
	p := new(Server)
	p.ctx = context

	mux := http.NewServeMux()
	p.AddRoutes(mux)

	p.Handler = mux
	return p
}

func (s *Server) AddRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/api/v1", handleBaseApi())
	mux.Handle("/api/v1/resume", handleGetResume())
}

func handleBaseApi() http.Handler {
	type baseResponse struct {
		Message string `json:"message"`
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			body := baseResponse{
				Message: "Hello, world!",
			}
			if err := encode(w, r, http.StatusOK, body); err != nil {
				fmt.Printf("error while encoding body %+v: %v\n", body, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		},
	)
}

func handleGetResume() http.Handler {
	type gitHubFile struct {
		Name        string `json:"name"`
		Path        string `json:"path"`
		Sha         string `json:"sha"`
		Size        string `json:"size"`
		Url         string `json:"url"`
		DownloadUrl string `json:"download_url"`
		Type        string `json:"type"`
		Content     string `json:"content"`
		Encoding    string `json:"encoding"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/jackjf28/resume/contents/jack_farrell_resume.pdf", nil)

		req.Header.Add("Accept", "application/vnd.github.object")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("PAT")))
		req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "failed to fetch resume", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		file, _ := decode[gitHubFile](resp)

		if file.Encoding != "base64" {
			http.Error(w, fmt.Sprintf("expected base64 encoding, got %s", file.Encoding), http.StatusInternalServerError)
			return
		}

		pdfData, _ := base64.StdEncoding.DecodeString(file.Content)
		// TODO: application/pdf -> directly serve via http?

		os.WriteFile("resume.pdf", pdfData, 0644)

		bytes, _ := io.ReadAll(resp.Body)
		w.Write(bytes)
	})
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Response) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
