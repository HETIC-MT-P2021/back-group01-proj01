package images

import (
	"fmt"
	cLog "image_gallery/logger"
	"image_gallery/router"
	"net/http"
)

// Handler is the home handler
type Handler struct {
	Logger *cLog.Logger
}

// Routes returns handler routes
func (h *Handler) Routes() router.Routes {
	return []router.Route{
		router.Route{
			Name:        "Get all images",
			Method:      "GET",
			Pattern:     "/images",
			HandlerFunc: h.getimages,
		},
	}
}

func (h *Handler) getimages(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "coucou, si si je suis une image")
}
