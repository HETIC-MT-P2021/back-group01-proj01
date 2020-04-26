package tag

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
			Name:        "Get a tag",
			Method:      "GET",
			Pattern:     "/tags/{id}",
			HandlerFunc: h.gettagByID,
		},
		router.Route{
			Name:        "Get all tags",
			Method:      "GET",
			Pattern:     "/tags",
			HandlerFunc: h.getAllTags,
		},
		router.Route{
			Name:        "Post a tag",
			Method:      "POST",
			Pattern:     "/tags",
			HandlerFunc: h.createTag,
		},
	}
}

func (h *Handler) gettagByID(w http.ResponseWriter, r *http.Request) {

	muxVars := mux.Vars(r)
	db := database.DbConn

	repository := Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	tag, err := repository.SelectTagByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve tag")
		return
	}

	if tag == nil {
		helpers.WriteJSON(w, http.StatusNotFound, tag)
		return
	}

	h.Logger.Infof("tag retrieved: %v", tag)
	helpers.WriteJSON(w, http.StatusOK, tag)
}

func (h *Handler) getAllTags(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn
	repository := Repository{Conn: db}

	tags, err := repository.retrieveAllTags()
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve tags")
		return
	}

	h.Logger.Infof("tags retrieved")
	helpers.WriteJSON(w, http.StatusOK, tags)
}

func (h *Handler) createTag(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn
	repository := Repository{Conn: db}

	var tag Tag

	err := helpers.ReadValidateJSON(w, r, &tag)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	err = repository.insertTag(&tag)

	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to save tag")
		return
	}

	h.Logger.Infof("saved tag: %v", tag)
	helpers.WriteJSON(w, http.StatusOK, tag)
}

func (h *Handler) updateTag(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	var tag Tag

	err = helpers.ReadValidateJSON(w, r, &tag)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	err = repository.updateTag(&tag, id)
	if err != nil {
		h.Logger.Error(err)
		return
	}
	h.Logger.Infof("updated tag: %v", tag)
	helpers.WriteJSON(w, http.StatusOK, tag)
}

func (h *Handler) deleteTag(w http.ResponseWriter, r *http.Request) {
	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}

	rowsAffected, err := repository.deleteTag(id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("%d tag deleted with ID: %v", rowsAffected, id)
	helpers.WriteJSON(w, http.StatusNoContent, "Category deleted")
}
