package home

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
			Name:        "Get all image",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: h.sayWelcome,
		},
	}
}

func (h *Handler) sayWelcome(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Good morning, nice day for fishin' ain't it? huhu")
}
