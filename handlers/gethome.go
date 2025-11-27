package handlers

import (
	"net/http"

	"github.com/jackjf28/resume-website/templates"
	"log/slog"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) GetHomePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("get home page request")
		name := "jack"
		c := templates.Home(name)
		err := templates.Layout(c, "My website", "/home").Render(r.Context(), w)

		if err != nil {
			slog.Error("error rendering template", "error", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})
}
