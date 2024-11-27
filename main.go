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

	// Obtener el puerto desde las variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor predeterminado
	}

	// Configurar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", handlers.UploadHandler)
	mux.HandleFunc("/convert", handlers.ConvertHandler)
	mux.HandleFunc("/download/", handlers.DownloadHandler)

	// Habilitar CORS
	handlerWithCORS := enableCORS(mux)

	log.Printf("Servidor corriendo en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
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
