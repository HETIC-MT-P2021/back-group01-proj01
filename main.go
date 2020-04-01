package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"image_gallery/images"
	logger "image_gallery/logger"
	"image_gallery/router"
	"net/http"
	"os"
	"strings"
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

	err := http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins(strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")),
			handlers.AllowedHeaders([]string{"Content-Type"}),
			handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE"}),
		)(muxRouter),
	)

	customLogger.Fatal(err)
}
