package main

import (
	"log"
	"net/http"
	"os"
	"wasabi/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargamos configuración (.env)
	// Esto permite que el puerto y las URLs sean dinámicos.
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, usando variables de sistema")
	}

	// 2. Instanciamos el ServeMux (Enrutador estándar)
	mux := http.NewServeMux()

	// 3. Mapeamos las rutas usando el módulo especializado
	// Esto mantiene este archivo Main corto y legible.
	handlers.MapRoutes(mux)

	// 4. Obtenemos el puerto de las variables de entorno o usamos 3000 por defecto
	port := os.Getenv("WASABI_PORT")
	if port == "" {
		port = "3000"
	}

	// 5. Configuración avanzada del servidor
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("🚀 Servidor Wasabi iniciado en http://localhost:%s", port)
	log.Printf("📡 Webhook listo en: http://localhost:%s/api/v1/wuzapi/webhook", port)

	// 6. Arrancamos el servidor
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error fatal al arrancar el servidor: ", err)
	}
}