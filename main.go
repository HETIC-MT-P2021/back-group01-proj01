package main

import (
	"image_gallery/images"
	logger "image_gallery/logger"
	"image_gallery/router"
	"net/http"
	"os"
)

func main() {
	customLogger := logger.GetLogger()

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	customLogger.Infof("Server started on port %s", port)
	apiRouter := router.Router{
		Logger: customLogger,
	}

	// Images handler
	apiRouter.AddHandler(&images.Handler{
		Logger: customLogger,
	})

	muxRouter := apiRouter.Configure()

	//TODO(athenais): add cors handling
	err := http.ListenAndServe(":8080", muxRouter)

	customLogger.Fatal(err)
}
