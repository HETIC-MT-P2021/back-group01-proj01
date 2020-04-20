package main

import (
	"image_gallery/category"
<<<<<<< HEAD
	"image_gallery/home"
=======
>>>>>>> ebf8e93948946df0f48926fba8d924877b4ab983
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"

	"image_gallery/database"
	"image_gallery/image"
	cLog "image_gallery/logger"
	"image_gallery/router"
)

func main() {
	logger := cLog.GetLogger()

	logger.Info("Server started on port 8080")
	apiRouter := router.Router{
		Logger: logger,
	}

	// Home handler
	apiRouter.AddHandler(&home.Handler{
		Logger: logger,
	})

	// Images handler
	apiRouter.AddHandler(&category.Handler{
		Logger: logger,
	})

	// Images handler
	apiRouter.AddHandler(&image.Handler{
<<<<<<< HEAD
=======
		Logger: logger,
	})

	// Category handler
	apiRouter.AddHandler(&category.Handler{
>>>>>>> ebf8e93948946df0f48926fba8d924877b4ab983
		Logger: logger,
	})

	err := database.Connect()

	if err != nil {
		logger.Fatalf("could not connect to db: %v", err)
	}

	muxRouter := apiRouter.Configure()

	err = http.ListenAndServe(
		":8080",
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins(strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
		)(muxRouter),
	)

	logger.Fatal(err)
}
