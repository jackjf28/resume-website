package handlers

import (
	"net/http"

	"github.com/jackjf28/resume-website/templates"
)

type GetResumeHandler struct{}

func NewGetResumeHandler() *GetResumeHandler {
	return &GetResumeHandler{}
}

func (h *GetResumeHandler) GetResumePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := templates.Resume()
		err := templates.Layout(c, "My website", "/resume").Render(r.Context(), w)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})
}
