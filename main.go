package main

import (
	"blank/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Crear directorios necesarios
	_ = os.MkdirAll("uploads", os.ModePerm)
	_ = os.MkdirAll("downloads", os.ModePerm)

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handlers.UploadHandler)
	mux.HandleFunc("/convert", handlers.ConvertHandler)
	mux.HandleFunc("/download/", handlers.DownloadHandler)

	// Habilitar CORS
	handlerWithCORS := enableCORS(mux)

	log.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCORS))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Responder rápido a las solicitudes OPTIONS (preflight)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
