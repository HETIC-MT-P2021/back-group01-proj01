<<<<<<< HEAD:app/home/handler.go
package home
=======
package image
>>>>>>> ebf8e93948946df0f48926fba8d924877b4ab983:app/image/handler.go

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
<<<<<<< HEAD:app/home/handler.go
			Pattern:     "/",
			HandlerFunc: h.sayWelcome,
=======
			Pattern:     "/image",
			HandlerFunc: h.getimages,
>>>>>>> ebf8e93948946df0f48926fba8d924877b4ab983:app/image/handler.go
		},
	}
}

func (h *Handler) sayWelcome(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Good morning, nice day for fishin' ain't it? huhu")
}
