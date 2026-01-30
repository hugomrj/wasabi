package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"wasabi/internal/wuzapi"
)

// Estructura exacta basada en tu log real
type WuzapiEvent struct {
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

	// 1. Parsear el formulario
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parseando formulario: %v", err)
		return
	}

	// 2. Obtener el jsonData
	rawJSON := r.FormValue("jsonData")
	var detail WuzapiEvent
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		log.Printf("❌ Error decodificando: %v", err)
		return
	}

	// 3. Capturar los datos reales (usamos Sender para responder)
	remitente := detail.Info.Sender
	texto := detail.Message.Conversation

	log.Printf("📦 Mensaje de %s: %s", remitente, texto)

	if detail.Info.IsFromMe {
		log.Println("⏭️ Ignorando mensaje propio")
		return
	}

	// 4. Responder
	if detail.Event == "Message" && remitente != "" {
		// ¡OJO AQUÍ!: Si el Token no viene en el Header, usa el que configuraste en el .env
		token := r.Header.Get("Token")
		if token == "" {
			token = "USER_TOKEN_1" // Forzamos el token si el webhook no lo trae
		}

		log.Printf("📩 Intentando responder a %s con Token: %s", remitente, token)
		
		err := wuzapi.SendMessage(token, remitente, "¡Hola Hugo! Tu bot Wasabi ya sabe leer y responder 🚀")
		if err != nil {
			log.Printf("❌ Error enviando: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}