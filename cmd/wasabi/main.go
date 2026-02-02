package main

import (
	"log"
	"net/http"
	"os"
	"wasabi/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno (.env)
	// No es fatal si no existe, por eso ignoramos el error
	godotenv.Load()

	port := os.Getenv("WASABI_PORT")
	if port == "" {
		port = "3000"
	}

	// 2. CREAR EL ROUTER (ServeMux)
	// Usamos un mux propio para soportar rutas dinámicas como /webhook/{instancia}
	mux := http.NewServeMux()

	// 3. REGISTRAR RUTAS DESDE EL HANDLER
	// Esto llama a tu archivo internal/handlers/routes.go y configura las rutas
	handlers.MapRoutes(mux)

	// 4. INICIAR EL SERVIDOR
	log.Printf("🚀 Wasabi levantado en el puerto %s", port)
	log.Printf("📌 Rutas activas: /health, /api/v1/health/ping, /webhook/{instancia}")

	// Importante: Pasamos 'mux' en lugar de 'nil'
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal("❌ Error al iniciar el servidor: ", err)
	}
}