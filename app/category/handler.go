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
			Pattern:     "/categories/{id}",
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
			Pattern:     "/categories",
			HandlerFunc: h.createCategory,
		},
		router.Route{
			Name:        "Update category",
			Method:      "PUT",
			Pattern:     "/categories/{id}",
			HandlerFunc: h.updateCategory,
		},
		router.Route{
			Name:        "Delete category",
			Method:      "DELETE",
			Pattern:     "/categories/{id}",
			HandlerFunc: h.deleteCategory,
		},
	}
}

func (h *Handler) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	muxVars := mux.Vars(r)
	db := database.DbConn

	repository := Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	category, err := repository.SelectCategoryByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve category")
		return
	}

	if category == nil {
		helpers.WriteJSON(w, http.StatusNotFound, category)
		return
	}

	h.Logger.Infof("category retrieved: %v", category)
	helpers.WriteJSON(w, http.StatusOK, category)
}

func (h *Handler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

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
	h.Logger.Infof("calling %v", r.URL.Path)

	db := database.DbConn
	repository := Repository{Conn: db}

	var category Category

	err := helpers.ReadValidateJSON(w, r, &category)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	err = repository.insertCategory(&category)

	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to save category")
		return
	}

	h.Logger.Infof("saved category: %v", category)
	helpers.WriteJSON(w, http.StatusOK, category)

}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

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
	err = repository.updateCategory(&category, id)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	h.Logger.Infof("updated category: %v", category)
	helpers.WriteJSON(w, http.StatusOK, category)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

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
