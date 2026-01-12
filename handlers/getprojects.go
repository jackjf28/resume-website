package handlers

import (
	"log/slog"
	"net/http"

	"github.com/jackjf28/resume-website/templates"
)

type ProjectHandler struct{}

var ilespum string = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed ultricies lectus ac quam ultricies viverra. Duis suscipit leo sed varius ultricies. Vivamus eget orci nec felis dignissim malesuada sit amet nec nisi. Nulla et tincidunt risus, vitae iaculis leo. Donec commodo leo ligula, non facilisis turpis maximus ac. Cras at maximus est. Nunc sapien ipsum, interdum ac fringilla eget, aliquam vitae tortor. Aliquam ut sapien sit amet ligula sollicitudin commodo auctor eget lorem. Vestibulum et finibus urna. Donec consectetur euismod sapien, pharetra interdum metus vestibulum nec. Etiam et ipsum in est rhoncus porttitor. Fusce non lectus aliquet, tristique quam id, efficitur quam. Nam quis orci elit.

Sed finibus purus tincidunt mi ultrices, sit amet ornare magna fringilla. Phasellus at blandit diam. Vestibulum sed quam finibus, dignissim turpis eu, congue ante. Cras id felis felis. Nunc fermentum posuere tortor non rutrum. Aliquam luctus nisl eu felis tempus, sit amet tempus ex mollis. Maecenas pellentesque ex sed augue dignissim, sed volutpat ante pretium.

Proin eget lacus mauris. Cras dolor est, cursus nec est non, efficitur tincidunt diam. Maecenas quis nunc euismod, varius dui nec, ornare mi. Ut sollicitudin nisl ac pellentesque tempor. Morbi at porta justo. Ut maximus malesuada lectus ut condimentum. Suspendisse suscipit eleifend facilisis. Suspendisse gravida in risus eu tincidunt. Donec sed commodo elit. Sed dignissim vulputate mi nec tincidunt. Duis in urna quis ex lacinia ultricies vel eget nulla.

Maecenas ac venenatis metus. Suspendisse tempus faucibus felis eu rutrum. Cras et erat ac lorem pharetra egestas. Ut tincidunt venenatis pretium. Aenean pharetra leo sed justo pulvinar pretium. Suspendisse iaculis nec mauris sed bibendum. Maecenas in leo eu augue scelerisque finibus. Nunc accumsan sagittis hendrerit. Vivamus ut dui quis ligula pharetra cursus. Nulla facilisi. Donec commodo arcu a mauris malesuada, quis posuere dui laoreet. Aenean nulla eros, faucibus ut dolor nec, tempor ornare diam.

Mauris vel dignissim sapien, ut lobortis nulla. Maecenas quis erat eu leo tristique volutpat. Suspendisse vel quam ipsum. Etiam sit amet felis magna. Etiam vitae nulla a nibh porttitor egestas. Maecenas ut ex velit. Mauris nulla sem, convallis a orci a, hendrerit gravida purus. `

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

func (p *ProjectHandler) GetHandler() http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		slog.Info("get projects request")
		
		projects := []templates.Project {
			{Name: "Project 1", Description: ilespum },
			{Name: "Project 2", Description: ilespum },
		}
		c := templates.Projects(projects)
		err := c.Render(r.Context(), w)
		if err != nil {
			slog.Error("error rendering template", "error", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	})
}
