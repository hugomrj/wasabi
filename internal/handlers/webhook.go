package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "wasabi/internal/models"
    "wasabi/internal/wuzapi"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        return
    }

    // Capturamos el Token del Header
    token := r.Header.Get("Token")

    var payload models.WuzapiRequest
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        log.Printf("Error: %v", err)
        return
    }

    // Respondemos solo si es un mensaje de texto
    if payload.Event == "Message" && payload.Data.Body != "" {
        log.Printf("Mensaje de %s (Instancia: %s)", payload.Data.From, token)
        
        // El "Hola" simple que pediste
        _ = wuzapi.SendMessage(token, payload.Data.From, "¡Hola! Recibido.")
    }

    w.WriteHeader(http.StatusOK)
}