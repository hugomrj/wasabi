package main

import (
    "log"
    "net/http"
    "os"
    "wasabi/internal/handlers"
    "github.com/joho/godotenv"
)

func main() {
    godotenv.Load() // Carga el .env

    port := os.Getenv("PORT")
    if port == "" {
        port = "8000"
    }

    http.HandleFunc("/webhook", handlers.WebhookHandler)

    log.Printf("🚀 Wasabi escuchando en puerto %s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}