package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "wasabi/internal/models"
    "wasabi/internal/wuzapi"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    // ESTO NOS DIRÁ SI ALGUIEN TOCÓ LA PUERTA
    log.Println("🔔 ¡Webhook invocado!")

    if r.Method != http.MethodPost {
        log.Println("❌ No es un POST")
        return
    }

    var payload models.WuzapiRequest
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        log.Printf("❌ Error decodificando: %v", err)
        return
    }

    log.Printf("📦 Payload recibido: %+v", payload)

    if payload.Data.FromMe {
        log.Println("⏭️ Ignorando mensaje propio (FromMe)")
        return
    }

    if payload.Event == "Message" {
        token := r.Header.Get("Token")
        log.Printf("📩 Respondiendo a: %s", payload.Data.From)
        _ = wuzapi.SendMessage(token, payload.Data.From, "¡Hola! Recibido.")
    }

    w.WriteHeader(http.StatusOK)
}	