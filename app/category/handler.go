package category

import (
	"github.com/gorilla/mux"
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
			Name:        "Get an image by id",
			Method:      "GET",
			Pattern:     "/category/{id}",
			HandlerFunc: h.getCategoryByID,
		},
		router.Route{
			Name:        "Get all categories",
			Method:      "GET",
			Pattern:     "/categories",
			HandlerFunc: h.getAllCategories,
		},
		router.Route{
			Name:        "Post category",
			Method:      "POST",
			Pattern:     "/category",
			HandlerFunc: h.createCategory,
		},
		router.Route{
			Name:        "Update category",
			Method:      "PUT",
			Pattern:     "/category/{id}",
			HandlerFunc: h.updateCategory,
		},
		router.Route{
			Name:        "Delete category",
			Method:      "DELETE",
			Pattern:     "/category/{id}",
			HandlerFunc: h.deleteCategory,
		},
	}
}

func (h *Handler) getCategoryByID(w http.ResponseWriter, r *http.Request) {

	muxVars := mux.Vars(r)
	db := database.DbConn

	repository := Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	category, err := repository.selectCategoryByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve category")
		return
	}

	if category == nil {
		helpers.WriteJSON(w, http.StatusNotFound, category)
	}

	h.Logger.Infof("category retrieved: %v", category)
	helpers.WriteJSON(w, http.StatusOK, category)
}

func (h *Handler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn
	repository := Repository{Conn: db}

	categories, err := repository.retrieveAllCategories()
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve categories")
		return
	}

	h.Logger.Infof("categories retrieved")
	helpers.WriteJSON(w, http.StatusOK, categories)
}

func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn
	repository := Repository{Conn: db}

	var category Category

	err := helpers.ReadValidateJSON(w, r, &category)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	categoryPosted, err := repository.insertCategory(&category)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to save category")
		return
	}

	h.Logger.Infof("saved category: %v", categoryPosted)
	helpers.WriteJSON(w, http.StatusOK, categoryPosted)
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	var category Category

	err = helpers.ReadValidateJSON(w, r, &category)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	categoryUpdated, err := repository.updateCategory(&category, id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("updated category: %v", categoryUpdated)
	helpers.WriteJSON(w, http.StatusOK, categoryUpdated)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	rowsAffected, err := repository.deleteCategory(id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("%d category deleted with ID: %v", rowsAffected, id)
	helpers.WriteJSON(w, http.StatusNoContent, "Category deleted")
}
