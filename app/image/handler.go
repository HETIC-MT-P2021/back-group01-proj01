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

	image, err := repository.selectImageByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve image")
		return
	}
	
	
	h.Logger.Printf("image: %v", image)

	if image == nil {
		helpers.WriteJSON(w, http.StatusNotFound, image)
		return
	}

	categoryRetrieved, err := categoryRepository.SelectCategoryByID(image.CategoryId)

	image.Category = categoryRetrieved
	
	h.Logger.Infof("image retrieved: %v", image)
	helpers.WriteJSON(w, http.StatusOK, image)
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
	
	var image Image
	err := helpers.ReadValidateJSON(w, r, &image)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to read image")
		return
	}
	
	err = imageRepository.insertImage(&image)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to save image")
		return
	}
	
	categoryRetrieved, err := categoryRepository.SelectCategoryByID(image.CategoryId)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, 
			"unable to retrieve image category")
		return
	}
	image.Category = categoryRetrieved

	h.Logger.Infof("saved image: %v", image)
	helpers.WriteJSON(w, http.StatusOK, image)
}

//func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
//	muxVars := mux.Vars(r)
//	id, err := helpers.ParseInt64(muxVars["id"])
//	if err != nil {
//		h.Logger.Error(err)
//		return
//	}
//	db := database.DbConn
//	repository := Repository{Conn: db}
//
//	var image Image
//
//	err = helpers.ReadValidateJSON(w, r, &image)
//	if err != nil {
//		h.Logger.Error(err)
//		return
//	}
//	
//	
//
//	categoryUpdated, err := repository.updateCategory(&image, id)
//	if err != nil {
//		h.Logger.Error(err)
//		return
//	}
//
//	h.Logger.Infof("updated category: %v", categoryUpdated)
//	helpers.WriteJSON(w, http.StatusOK, categoryUpdated)
//}
