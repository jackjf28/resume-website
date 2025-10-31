package server

import (
	"context"
	"net/http"
	"os"

	"github.com/jackjf28/resume-website/github"
	"github.com/jackjf28/resume-website/handlers"
	"github.com/jackjf28/resume-website/middleware"
	"github.com/jackjf28/resume-website/services"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	http.Handler
	ctx context.Context
	resumeHandler *handlers.ResumeHandler
	homeHandler *handlers.HomeHandler
}

func NewServer(context context.Context) http.Handler {

	githubClient := github.NewClient(os.Getenv("PAT"))
	resumeService := services.NewResumeService(githubClient)
	resumeHandler := handlers.NewResumeHandler(resumeService)

	homeHandler := handlers.NewHomeHandler()

	p := &Server{
		ctx: context,
		resumeHandler: resumeHandler,
		homeHandler: homeHandler,
	}
	mux := http.NewServeMux()


	p.AddRoutes(mux)
	p.Handler = middleware.TextHTMLMiddleware(
		middleware.CSPMiddleware(mux),
	)
	return p
}

func (s *Server) AddRoutes(mux *http.ServeMux) {
	// serve static content
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/", http.NotFoundHandler())
	mux.Handle("/home", s.homeHandler.GetHomePage())
	mux.Handle("/api/v1", http.NotFoundHandler())
	mux.Handle("/api/v1/resume", s.resumeHandler.GetResume())
}


//func handleBaseAPI() http.Handler {
//	type baseResponse struct {
//		Message string `json:"message"`
//	}
//	return http.HandlerFunc(
//		func(w http.ResponseWriter, r *http.Request) {
//			body := baseResponse{
//				Message: "Hello, world!",
//			}
//			if err := utils.Encode(w, r, http.StatusOK, body); err != nil {
//				fmt.Printf("error while encoding body %+v: %v\n", body, err)
//				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//			}
//		},
//	)
//}
