package Module

import (
	"net/http"
	"strings"
	"fmt"
	"io"
	"os"
	"encoding/json"
)

func UploadFiles(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	f,errorCopy := os.OpenFile("./"+ header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if errorCopy != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}
	defer f.Close()

	io.Copy(f, file)

	return
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
