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

    token := r.Header.Get("Token")

    var payload models.WuzapiRequest
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        log.Printf("Error: %v", err)
        return
    }

    // 1. EL FILTRO: Si FromMe es true, salimos silenciosamente
    if payload.Data.FromMe {
        return 
    }

    // 2. LA LÓGICA: Solo si es mensaje y no es mío
    if payload.Event == "Message" && payload.Data.Body != "" {
        log.Printf("📩 Mensaje de %s (Token: %s)", payload.Data.From, token)
        
        // Enviamos la respuesta
        err := wuzapi.SendMessage(token, payload.Data.From, "¡Hola! Recibido.")
        if err != nil {
            log.Printf("❌ Error al enviar: %v", err)
        }
    }

    w.WriteHeader(http.StatusOK)
}