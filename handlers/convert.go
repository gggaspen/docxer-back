package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	// Decodificar el nombre del archivo
	fileName, err := url.QueryUnescape(req.FileName)
	if err != nil {
		http.Error(w, "Nombre de archivo inválido", http.StatusBadRequest)
		return
	}

	inputPath := filepath.Join("uploads", fileName)
	outputPath := filepath.Join("downloads", fileName)

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

	w.Write([]byte(fmt.Sprintf("Archivo convertido: %s", fileName)))
}

// DownloadHandler permite descargar el archivo convertido
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el nombre del archivo desde la URL
	fileName, err := url.QueryUnescape(filepath.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "Nombre de archivo inválido", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("downloads", fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
