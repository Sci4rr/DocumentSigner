package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logMessage("No .env file found")
	}
}

type Document struct {
	ID         string
	FileName   string
	Owner      string
	SharedWith []string
	FilePath   string
}

func logMessage(message string) {
	fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logMessage("Error parsing multipart form: " + err.Error())
		return
	}

	uploadFile, fileHeader, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logMessage("Error retrieving the file: " + err.Error())
		return
	}
	defer uploadFile.Close()

	destFilePath := filepath.Join(os.Getenv("DOCUMENT_STORAGE_PATH"), fileHeader.Filename)
	destFile, err := os.Create(destFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logMessage("Error creating the file: " + err.Error())
		return
	}
	defer destFile.Close()

	if _, err = io.Copy(destFile, uploadFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logMessage("Error copying the file: " + err.Error())
		return
	}

	logMessage("File uploaded successfully: " + fileHeader.Filename)
	fmt.Fprintf(w, "File uploaded successfully: %s", fileHeader.Filename)
}

func handleRetrieve(w http.ResponseWriter, r *http.Request) {
	requestedFileName := r.URL.Query().Get("filename")
	requestedFilePath := filepath.Join(os.Getenv("DOCUMENT_STORAGE_PATH"), requestedFileName)

	if _, err := os.Stat(requestedFilePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		logMessage("File not found: " + requestedFileName)
		return
	}

	logMessage("File retrieved successfully: " + requestedFileName)
	http.ServeFile(w, r, requestedFilePath)
}

func handleUpdatePermissions(w http.ResponseWriter, r *http.Request) {
	// Without specific details on permission updating, this is a placeholder function
	logMessage("Document permissions updated")
	fmt.Fprintf(w, "Document permissions updated")
}

func setupServerRoutes() {
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/retrieve", handleRetrieve)
	http.HandleFunc("/permissions", handleUpdatePermissions)
}

func main() {
	setupServerRoutes()
	logMessage("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}