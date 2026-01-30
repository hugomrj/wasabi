package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wasabi/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Intentamos cargar el .env (en el servidor lo leerá de /srv/wasabi/.env)
	godotenv.Load()

	port := os.Getenv("WASABI_PORT")
	if port == "" {
		port = "3000" // Cambié a 3000 por defecto para que coincida con tu config
	}

	// --- RUTAS ---
	
	// Ruta de Salud (Para verificar que el binario corre)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "🌿 Wasabi está en línea (Puerto %s)", port)
	})

	// Ruta del Webhook para Wuzapi
	http.HandleFunc("/webhook", handlers.WebhookHandler)

	// --- INICIO ---
	log.Printf("🚀 Wasabi levantado en el puerto %s", port)
	
	// Escuchar en todas las interfaces para que sea accesible desde afuera
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("❌ Error al iniciar el servidor: ", err)
	}
}