package main

import (
	"image_gallery/category"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"

	"image_gallery/database"
	"image_gallery/images"
	cLog "image_gallery/logger"
	"image_gallery/router"
)

func main() {
	logger := cLog.GetLogger()

	logger.Info("Server started on port 8080")
	apiRouter := router.Router{
		Logger: logger,
	}

	// Images handler
	apiRouter.AddHandler(&images.Handler{
		Logger: logger,
	})
	apiRouter.AddHandler(&category.Handler{
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
