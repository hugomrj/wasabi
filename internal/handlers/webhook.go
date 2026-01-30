package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"wasabi/internal/wuzapi"
)

// Estructura corregida para seguir el mapa real del log
type WuzapiResponse struct {
	EventData struct {
		Info struct {
			Sender   string `json:"Sender"`
			SenderAlt string `json:"SenderAlt"` // Usamos este si el primero falla
			IsFromMe bool   `json:"IsFromMe"`
		} `json:"Info"`
		Message struct {
			Conversation string `json:"conversation"`
		} `json:"Message"`
		Type string `json:"type"`
	} `json:"event"` // IMPORTANTE: Todo está dentro de "event"
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔔 ¡Webhook invocado!")

	r.ParseForm()
	rawJSON := r.FormValue("jsonData")
	
	var wrapper WuzapiResponse
	if err := json.Unmarshal([]byte(rawJSON), &wrapper); err != nil {
		log.Printf("❌ Error decodificando: %v", err)
		return
	}

	// Extraemos los datos del nivel "event"
	info := wrapper.EventData.Info
	msg := wrapper.EventData.Message

	// Elegimos el remitente: Si Sender tiene ":", preferimos SenderAlt que es el número limpio
	remitente := info.Sender
	if info.SenderAlt != "" {
		remitente = info.SenderAlt
	}
	// Limpiamos el @s.whatsapp.net si existe
	remitente = strings.Split(remitente, "@")[0]
	remitente = strings.Split(remitente, ":")[0]

	log.Printf("📦 Mensaje de [%s]: %s", remitente, msg.Conversation)

	if info.IsFromMe {
		log.Println("⏭️ Ignorando mensaje propio")
		return
	}

	if remitente != "" && msg.Conversation != "" {
		// Usamos el token fijo que sabemos que funciona
		token := "USER_TOKEN_1" 
		
		log.Printf("📩 Respondiendo a %s...", remitente)
		err := wuzapi.SendMessage(token, remitente, "¡Hugo! Ahora sí capturé tu mensaje: " + msg.Conversation)
		if err != nil {
			log.Printf("❌ Error enviando: %v", err)
		} else {
			log.Println("✅ Respuesta enviada con éxito")
		}
	}

	w.WriteHeader(http.StatusOK)
}