package handlers

import (
    "io"
    "log"
    "net/http"
    // "encoding/json"       <-- Comentado para que compile
    // "wasabi/internal/wuzapi" <-- Comentado para que compile
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("🔔 ¡Webhook invocado!")

    // Leemos TODO lo que venga en el cuerpo
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("❌ Error leyendo cuerpo: %v", err)
        return
    }

    // ESTO ES LO QUE NECESITAMOS VER
    log.Printf("DEBUG - CUERPO CRUDO: %s", string(body))
    log.Printf("DEBUG - CONTENT-TYPE: %s", r.Header.Get("Content-Type"))

    w.WriteHeader(http.StatusOK)
}