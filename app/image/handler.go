package image

import (
	"fmt"
	"github.com/gorilla/mux"
	"image_gallery/category"
	"image_gallery/database"
	"image_gallery/helpers"
	cLog "image_gallery/logger"
	"image_gallery/router"
	"image_gallery/tag"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
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
		router.Route{
			Name:        "Upload an image",
			Method:      "POST",
			Pattern:     "/upload/{id}",
			HandlerFunc: h.upload,
		},
	}
}

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
// UploadPath const to set upload path for all images
const UploadPath = "/go/uploads/"

func (h *Handler) getImagebyID(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	muxVars := mux.Vars(r)
	db := database.DbConn

	repository := Repository{Conn: db}
	categoryRepository := category.Repository{Conn: db}
	tagRepository := tag.Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	imageSelected, err := repository.selectImageByID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve image")
		return
	}

	if imageSelected == nil {
		helpers.WriteJSON(w, http.StatusNotFound, imageSelected)
		return
	}

	categoryRetrieved, err := categoryRepository.SelectCategoryByID(imageSelected.CategoryID)

	imageSelected.Category = categoryRetrieved

	tags, err := tagRepository.GetAllTagsByImageID(id)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve tags")
		return
	}

	imageSelected.TagsNames = tags

	h.Logger.Infof("image retrieved: %v", imageSelected)
	helpers.WriteJSON(w, http.StatusOK, imageSelected)
}

func (h *Handler) getAllImages(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	db := database.DbConn
	repository := Repository{Conn: db}

	filters := make(map[filterName]interface{})

	order := r.URL.Query().Get(string(filterByDateOfUpdate))
	if order != "" {
		filters[filterByDateOfUpdate] = order
	}

	tagID, _ := helpers.ParseInt64(r.URL.Query().Get(string(filterByTag)))

	if tagID != 0 {
		filters[filterByTag] = tagID
	}

	categoryID, _ := helpers.ParseInt64(r.URL.Query().Get(string(filterByCategory)))

	if categoryID != 0 {
		filters[filterByCategory] = categoryID
	}

	images, err := repository.retrieveAllImages(filters, db)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve images")
		return
	}

	if images == nil {
		helpers.WriteJSON(w, http.StatusNotFound, "no images found")
		return
	}

	h.Logger.Infof("images retrieved")
	helpers.WriteJSON(w, http.StatusOK, images)
}

func (h *Handler) createImage(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	db := database.DbConn
	imageRepository := Repository{Conn: db}
	categoryRepository := category.Repository{Conn: db}
	tagRepository := tag.Repository{Conn: db}

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

	if imageToCreate.TagsNames != nil {
		err = saveTags(imageRepository, tagRepository, &imageToCreate)
		if err != nil {
			h.Logger.Error(err)
			helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not save tags linked to image")
		}
	}

	categoryRetrieved, err := categoryRepository.SelectCategoryByID(imageToCreate.CategoryID)
	if err != nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "unable to retrieve image category")
		return
	}
	imageToCreate.Category = categoryRetrieved

	h.Logger.Infof("saved image: %v", imageToCreate)
	helpers.WriteJSON(w, http.StatusOK, imageToCreate)
}

func (h *Handler) updateImage(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	muxVars := mux.Vars(r)
	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}
	db := database.DbConn
	repository := Repository{Conn: db}
	tagRepository := tag.Repository{Conn: db}

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

	if image.TagsNames != nil {
		err = saveTags(repository, tagRepository, &image)
		if err != nil {
			h.Logger.Error(err)
			helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not save tags linked to image")
		}
	}

	h.Logger.Infof("updated image: %v", image)
	helpers.WriteJSON(w, http.StatusOK, image)
}

func (h *Handler) deleteImage(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	muxVars := mux.Vars(r)

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	db := database.DbConn
	repository := Repository{Conn: db}

	image, err := repository.selectImageByID(id)
	if err != nil {
		h.Logger.Error(err)
		return
	}

	if image == nil {
		h.Logger.Error(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "this image does not exist")
		return
	}

	// Hard delete mode deletes both image and image metadata
	if r.URL.Query().Get("delete_mode") == "hard" {

		rowsAffected, err := repository.deleteImage(id)
		if err != nil {
			h.Logger.Error(err)
			return
		}

		h.Logger.Infof("%d image deleted with ID: %v", rowsAffected, id)
	}

	path := UploadPath + muxVars["id"] + "/" + image.Slug + image.Type

	if _, err := os.Stat(path); os.IsExist(err) {
		err = os.Remove(path)
		if err != nil {
			h.Logger.Error(err)
			helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not delete image")
			return
		}
	}

	h.Logger.Infof("image deleted")
	helpers.WriteJSON(w, http.StatusNoContent, "Image deleted")

}

func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	h.Logger.Infof("calling %v", r.URL.Path)

	muxVars := mux.Vars(r)
	db := database.DbConn

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	repository := Repository{Conn: db}

	id, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		h.Logger.Error(err)
		return
	}

	image, err := repository.selectImageByID(id)
	if err != nil {
		h.Logger.Errorf("could not retrieve image by id : %v", err)
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "Could not check if image has already been uploaded")
	}

	if image == nil {
		h.Logger.Infof("image metadata with id %d does not exists", id)
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "image id"+muxVars["id"]+" does not exist")
		return
	}

	if image.Type != "" {
		h.Logger.Errorf("image has already been uploaded to file server")
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "You already have uploaded this image")
		return
	}

	file, handle, err := r.FormFile("file")
	if err != nil {
		h.Logger.Errorf("get file from form data failed : %v", err)
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "Could not upload file")
	}
	defer file.Close()

	fileSize := handle.Size
	// validate file size
	if fileSize > maxUploadSize {
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "File cannot exceed 2MB")
		return
	}

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		err = saveFile(w, file, handle, image)
		if err != nil {
			h.Logger.Errorf("could not save file: %v", err)
			helpers.WriteErrorJSON(w, http.StatusBadRequest, "File could not be uploaded")
			return
		}
	case "image/png":
		err = saveFile(w, file, handle, image)
		if err != nil {
			h.Logger.Errorf("could not save file: %v", err)
			helpers.WriteErrorJSON(w, http.StatusBadRequest, "File could not be uploaded")
			return
		}
	default:
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "The format file is not valid.")
	}
}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader, image *Image) error {

	data, err := ioutil.ReadAll(file)
	db := database.DbConn
	repository := Repository{Conn: db}
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	dirName := strconv.FormatInt(image.ID, 10)
	fileName := image.Slug
	extensions, err := mime.ExtensionsByType(handle.Header.Get("Content-Type"))

	if _, err := os.Stat(UploadPath + dirName); os.IsNotExist(err) {
		err = os.Mkdir(UploadPath+dirName, 0755)
		if err != nil {
			return fmt.Errorf("could not write directory: %v", err)
		}
	}

	err = ioutil.WriteFile(UploadPath+dirName+"/"+fileName+extensions[0], data, 0755)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	image.Type = extensions[0]

	err = repository.updateImage(image, image.ID)
	if err != nil {
		return fmt.Errorf("could not update image type: %v", err)
	}

	helpers.WriteJSON(w, http.StatusCreated, "File uploaded successfully!.")
	return nil
}

func saveTags(imageRepository Repository, tagRepository tag.Repository, imageTagged *Image) error {
	for _, tagName := range imageTagged.TagsNames {

		tagByName, err := tagRepository.SelectTagBy("name", tagName)
		if err != nil {
			return fmt.Errorf("could not check if tag already exists %v", err)
		}

		if tagByName != nil {
			err = imageRepository.linkTagToImage(imageTagged.ID, tagByName.ID)
			if err != nil {
				return fmt.Errorf("could not save tag %v", err)
			}
			continue
		}

		newTag := &tag.Tag{Name: tagName}

		err = tagRepository.InsertTag(newTag)
		if err != nil {
			return fmt.Errorf("could not save tag %v", err)
		}
		err = imageRepository.linkTagToImage(imageTagged.ID, newTag.ID)
		if err != nil {
			return fmt.Errorf("could not save tag %v", err)
		}
	}
	return nil
}
