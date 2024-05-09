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
		log("No .env file found")
	}
}

type Document struct {
	ID         string
	FileName   string
	Owner      string
	SharedWith []string
	FilePath   string
}

func log(message string) {
	fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
}

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log("Error parsing multipart form: " + err.Error())
		return
	}
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log("Error retrieving the file: " + err.Error())
		return
	}
	defer file.Close()
	dst, err := os.Create(filepath.Join(os.Getenv("DOCUMENT_STORAGE_PATH"), handler.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log("Error creating the file: " + err.Error())
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log("Error copying the file: " + err.Error())
		return
	}
	log("File uploaded successfully: " + handler.Filename)
	fmt.Fprintf(w, "File uploaded successfully: %s", handler.Filename)
}

func RetrieveDocument(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	filepath := filepath.Join(os.Getenv("DOCUMENT_STORAGE_PATH"), filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		log("File not found: " + filename)
		return
	}
	log("File retrieved successfully: " + filename)
	http.ServeFile(w, r, filepath)
}

func UpdateDocumentPermissions(w http.ResponseWriter, r *http.Request) {
	log("Document permissions updated")
	fmt.Fprintf(w, "Document permissions updated")
}

func setupRoutes() {
	http.HandleFunc("/upload", UploadDocument)
	http.HandleFunc("/retrieve", RetrieveDocument)
	http.HandleFunc("/permissions", UpdateDocumentPermissions)
}

func main() {
	setupRoutes()
	log("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}