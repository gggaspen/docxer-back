package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nguyenthenguyen/docx"
)

// UploadHandler maneja la carga de archivos
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Método no permitido")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error al leer el archivo")
		http.Error(w, "Error al leer el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join("uploads", header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error al guardar el archivo")
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println("Error al guardar el archivo")
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Archivo subido exitosamente: %s", header.Filename)))
}

// ConvertHandler maneja la conversión de archivos
func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	type Conversion struct {
		FileName string            `json:"file_name"`
		Changes  map[string]string `json:"changes"`
	}

	var req Conversion
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	inputPath := filepath.Join("uploads", req.FileName)
	outputPath := filepath.Join("downloads", req.FileName)

	doc, err := docx.ReadDocxFile(inputPath)
	if err != nil {
		http.Error(w, "Error al leer el archivo", http.StatusInternalServerError)
		return
	}
	defer doc.Close()

	content := doc.Editable()

	for oldWord, newWord := range req.Changes {
		content.Replace(fmt.Sprintf(" %s ", oldWord), fmt.Sprintf(" %s ", newWord), -1)
	}

	err = content.WriteToFile(outputPath)
	if err != nil {
		http.Error(w, "Error al guardar el archivo convertido", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Archivo convertido: %s", req.FileName)))
}

// DownloadHandler permite descargar el archivo convertido
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := filepath.Base(r.URL.Path)
	filePath := filepath.Join("downloads", fileName)

	http.ServeFile(w, r, filePath)
}
