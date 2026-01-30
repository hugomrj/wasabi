package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"wasabi/internal/wuzapi"
)

// Estructura temporal para capturar lo que manda Wuzapi (Asternic version)
type RawWuzapi struct {
	Event    string `json:"event"`
	Instance string `json:"instanceName"`
	JSONData string `json:"jsonData"` // Aquí viene el mensaje real como texto
}

// Estructura para el contenido de jsonData
type MessageDetail struct {
	Event string `json:"type"`
	Info  struct {
		Sender   string `json:"Sender"`
		IsFromMe bool   `json:"IsFromMe"`
	} `json:"Info"`
	Message struct {
		Conversation string `json:"conversation"`
	} `json:"Message"`
}




func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("🔔 ¡Webhook invocado!")

    // Leemos TODO lo que venga en el cuerpo, sea lo que sea
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("❌ Error leyendo cuerpo: %v", err)
        return
    }

    // IMPRIMIMOS EL CUERPO TAL CUAL LLEGA
    log.Printf("DEBUG - CUERPO CRUDO: %s", string(body))
    log.Printf("DEBUG - CONTENT-TYPE: %s", r.Header.Get("Content-Type"))

    // Por ahora solo respondemos OK para que Wuzapi no reintente
    w.WriteHeader(http.StatusOK)
}