package image

import (
	"github.com/gorilla/mux"
	"image_gallery/category"
	"image_gallery/database"
	"image_gallery/helpers"
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
			Pattern:     "/images/{id}",
			HandlerFunc: h.getImagebyID,
		},
		router.Route{
			Name:        "Get all images",
			Method:      "GET",
			Pattern:     "/images",
			HandlerFunc: h.getAllImages,
		},
		router.Route{
			Name:        "Post an image",
			Method:      "POST",
			Pattern:     "/images",
			HandlerFunc: h.createImage,
		},
		router.Route{
			Name:        "Update an image",
			Method:      "PUT",
			Pattern:     "/images/{id}",
			HandlerFunc: h.updateImage,
		},
		router.Route{
			Name:        "Delete an image",
			Method:      "DELETE",
			Pattern:     "/images/{id}",
			HandlerFunc: h.deleteImage,
		},
	}
}

func (h *Handler) getImagebyID(w http.ResponseWriter, r *http.Request) {

	muxVars := mux.Vars(r)
	db := database.DbConn

	repository := Repository{Conn: db}
	categoryRepository := category.Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Printf("VAR %v", id)

	imageSelected, err := repository.selectImageByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve image")
		return
	}

	h.Logger.Printf("image: %v", imageSelected)

	if imageSelected == nil {
		helpers.WriteJSON(w, http.StatusNotFound, imageSelected)
		return
	}

	categoryRetrieved, err := categoryRepository.SelectCategoryByID(imageSelected.CategoryID)

	imageSelected.Category = categoryRetrieved

	h.Logger.Infof("image retrieved: %v", imageSelected)
	helpers.WriteJSON(w, http.StatusOK, imageSelected)
}

func (h *Handler) getAllImages(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn
	repository := Repository{Conn: db}

	images, err := repository.retrieveAllImages()
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve images")
		return
	}

	h.Logger.Infof("images retrieved")
	helpers.WriteJSON(w, http.StatusOK, images)
}

func (h *Handler) createImage(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn
	imageRepository := Repository{Conn: db}
	categoryRepository := category.Repository{Conn: db}
	h.Logger.Debugf("calling %v", r.URL.Path)

	var imageToCreate Image
	err := helpers.ReadValidateJSON(w, r, &imageToCreate)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to read image")
		return
	}

	err = imageRepository.insertImage(&imageToCreate)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to save image")
		return
	}

	categoryRetrieved, err := categoryRepository.SelectCategoryByID(imageToCreate.CategoryID)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError,
			"unable to retrieve image category")
		return
	}
	imageToCreate.Category = categoryRetrieved

	h.Logger.Infof("saved image: %v", imageToCreate)
	helpers.WriteJSON(w, http.StatusOK, imageToCreate)
}

func (h *Handler) updateImage(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	var image Image

	err = helpers.ReadValidateJSON(w, r, &image)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	err = repository.updateImage(&image, id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("updated image: %v", image)
	helpers.WriteJSON(w, http.StatusOK, image)
}

func (h *Handler) deleteImage(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	rowsAffected, err := repository.deleteImage(id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("%d image deleted with ID: %v", rowsAffected, id)
	helpers.WriteJSON(w, http.StatusNoContent, "Image deleted")
}
