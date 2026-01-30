package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"wasabi/internal/wuzapi"
)


// Esta estructura coincide con el JSON que viene dentro de jsonData
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

	// 1. Parseamos el formulario (porque viene como x-www-form-urlencoded)
	if err := r.ParseForm(); err != nil {
		log.Printf("❌ Error parseando formulario: %v", err)
		return
	}

	// 2. Extraemos el campo 'jsonData'
	rawJSON := r.FormValue("jsonData")
	if rawJSON == "" {
		log.Println("⚠️ No se encontró jsonData en la petición")
		return
	}

	// 3. Decodificamos el JSON que estaba escondido en el formulario
	var detail WuzapiEvent
	if err := json.Unmarshal([]byte(rawJSON), &detail); err != nil {
		log.Printf("❌ Error decodificando JSON interno: %v", err)
		return
	}

	log.Printf("📦 Mensaje de %s: %s", detail.Info.Sender, detail.Message.Conversation)

	// 4. Filtro: No respondernos a nosotros mismos
	if detail.Info.IsFromMe {
		log.Println("⏭️ Ignorando mensaje propio")
		return
	}

	// 5. Responder si es un mensaje
	if detail.Event == "Message" {
		token := r.Header.Get("Token")
		log.Printf("📩 Enviando respuesta a: %s", detail.Info.Sender)
		
		// Enviamos la respuesta usando tu módulo wuzapi
		err := wuzapi.SendMessage(token, detail.Info.Sender, "¡Recibido! Wasabi está funcionando 🚀")
		if err != nil {
			log.Printf("❌ Error enviando: %v", err)
		}
	}

	w.WriteHeader(http.StatusOK)
}