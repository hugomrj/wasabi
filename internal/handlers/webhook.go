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

	// 1. Leer el cuerpo crudo para que no de error de decodificación
	body, _ := io.ReadAll(r.Body)
	
	var raw RawWuzapi
	if err := json.Unmarshal(body, &raw); err != nil {
		log.Printf("❌ Error decodificando Raw: %v", err)
		return
	}

	// 2. Decodificar el jsonData que viene dentro
	var detail MessageDetail
	if err := json.Unmarshal([]byte(raw.JSONData), &detail); err != nil {
		log.Printf("❌ Error decodificando jsonData: %v", err)
		return
	}

	log.Printf("📦 Mensaje recibido de %s: %s", detail.Info.Sender, detail.Message.Conversation)

	// 3. Filtros
	if detail.Info.IsFromMe {
		log.Println("⏭️ Ignorando mensaje propio")
		return
	}

	if detail.Event == "Message" {
		token := r.Header.Get("Token")
		// Limpiar el Sender (Wuzapi manda 595...:44@s.whatsapp.net, necesitamos solo el número)
		// Aquí lo mandamos tal cual, Wuzapi suele limpiarlo solo al enviar
		log.Printf("📩 Respondiendo a: %s", detail.Info.Sender)
		err := wuzapi.SendMessage(token, detail.Info.Sender, "¡Hola! Recibido desde Wasabi.")
		if err != nil {
			log.Printf("❌ Error al enviar respuesta: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}