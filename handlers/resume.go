package handlers

import (
	"fmt"
	"net/http"

	"github.com/jackjf28/resume-website/services"
)

const pdfContentType = "application/pdf"

type ResumeHandler struct {
	resumeService *services.ResumeService
}

func NewResumeHandler(resumeService *services.ResumeService) *ResumeHandler {
	return &ResumeHandler{
		resumeService: resumeService,
	}
}

func (h *ResumeHandler) GetResume() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		pdfData, err := h.resumeService.GetResumePDF(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch resume: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", pdfContentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))
		w.Header().Set("Content-Disposition", "inline; filename=\"resume.pdf\"")

		w.Write(pdfData)
	})
}
