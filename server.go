package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const jsonContentType string = "application/json"

type Server struct {
	http.Handler
}

func NewServer(context context.Context) http.Handler {
	p := new(Server)

	mux := http.NewServeMux()
	addRoutes(mux)

	p.Handler = mux
	return p
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/api/v1", handleBaseApi())
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

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}
