package interfaces

import "net/http"

type TemplateHandler interface {
	GetHandler() http.Handler
}
