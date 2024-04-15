package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/common"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		HandleError(w, r, common.UNABLE_TO_FIND_RESOURCE, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(30 << 20)
	file, _, err := r.FormFile("image")

	if err != nil {
		HandleError(w, r, common.UNABLE_TO_READ, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Determine the MIME type of the uploaded file
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		HandleError(w, r, common.UNABLE_TO_READ, "Error reading the file", http.StatusBadRequest)
		return
	}
	contentType := http.DetectContentType(buffer)

	// Map the MIME type to the corresponding file extension
	var fileExtension string
	switch contentType {
	case "image/jpeg":
		fileExtension = ".jpg"
	case "image/png":
		fileExtension = ".png"
	default:
		HandleError(w, r, common.UNABLE_TO_READ, "Invalid file type. Only JPEG and PNG are allowed", http.StatusBadRequest)
		return
	}

	// Reset the file reader to the beginning
	file.Seek(0, io.SeekStart)

	tempDir := "images"
	os.MkdirAll(tempDir, 0755)

	tempFile, err := os.CreateTemp(tempDir, fmt.Sprintf("user-upload-*%s", fileExtension))

	if err != nil {
		HandleError(w, r, common.UNABLE_TO_SAVE, "Error creating temporary file", http.StatusInternalServerError)
		return
	}

	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)

	if err != nil {
		HandleError(w, r, common.UNABLE_TO_SAVE, "Error copying file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully uploaded file: %s", filepath.Base(tempFile.Name()))
}

func ServeImage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		HandleError(w, r, common.UNABLE_TO_FIND_RESOURCE, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	imageName := r.URL.Path[len("/v1/images/"):]
	imagePath := filepath.Join("images", imageName)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		HandleError(w, r, common.UNABLE_TO_FIND_RESOURCE, "Image not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, imagePath)
}
