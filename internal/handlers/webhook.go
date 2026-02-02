package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time" 
	"wasabi/internal/wuzapi"
)

// WebhookPayload representa la estructura genérica del mensaje entrante
type WebhookPayload struct {
	EventData struct {
		Info struct {
			Sender    string `json:"Sender"`
			SenderAlt string `json:"SenderAlt"`
			IsFromMe  bool   `json:"IsFromMe"`
		} `json:"Info"`
		Message struct {
			Conversation string `json:"conversation"`
			ExtendedText struct {
				Text string `json:"text"`
			} `json:"extendedTextMessage"`
		} `json:"Message"`
		Type string `json:"type"`
	} `json:"event"`
}

// 1. SEMÁFORO GLOBAL: Solo permite 1 consulta a la vez a la IA.
// Esto garantiza que los mensajes hagan fila y no saturen la RAM.
var iaSemaphore = make(chan struct{}, 1)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	// --- LOG GIGANTE PARA EL TOKEN ---
	log.Println("################################################")
	log.Println("--- INICIO HEADERS ---")
    for nombre, valores := range r.Header {
        log.Printf("Header: %s = %v", nombre, valores)
    }
    log.Println("--- FIN HEADERS ---")
	log.Println("################################################")

	r.ParseForm()
	rawJSON := r.FormValue("jsonData")
	if rawJSON == "" {
		return
	}
	// --- ESTO TE MOSTRARÁ TODO EN LA CONSOLA ---
	if rawJSON != "" {
		log.Printf("📥 JSON RECIBIDO: %s", rawJSON)
	} else {
		log.Printf("⚠️ Webhook llamado pero 'jsonData' está vacío")
		return
	}
	// -------------------------------------------



	go func(data string) {
		var payload WebhookPayload
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("❌ Error decodificando payload: %v", err)
			return
		}

		if payload.EventData.Info.IsFromMe {
			return
		}

		prompt := payload.EventData.Message.Conversation
		if prompt == "" {
			prompt = payload.EventData.Message.ExtendedText.Text
		}

		if prompt == "" {
			return
		}

		remitente := payload.EventData.Info.Sender
		if payload.EventData.Info.SenderAlt != "" {
			remitente = payload.EventData.Info.SenderAlt
		}
		remitente = strings.Split(strings.Split(remitente, "@")[0], ":")[0]

		// 2. LLAMADA AL DOMINIO: GetExternalResponse ahora gestiona la fila
		respuestaIA := GetExternalResponse(prompt)

		token := "USER_TOKEN_1"
		err := wuzapi.SendMessage(token, remitente, respuestaIA)
		if err != nil {
			log.Printf("❌ Error enviando a %s: %v", remitente, err)
		} else {
			log.Printf("✅ Respuesta enviada con éxito a %s", remitente)
		}
	}(rawJSON)
}





func GetExternalResponse(prompt string) string {
	// 3. ENTRADA A LA FILA: Si hay otro proceso, este espera aquí.
	log.Printf("⏳ Mensaje en espera de turno para IA...")
	iaSemaphore <- struct{}{}
	defer func() { <-iaSemaphore }() // Libera el turno al salir

	const targetURL = "https://japo.click/charlette/ask"
	const maxRetries = 2

	log.Printf("📩 Procesando IA ahora...")

	for i := 0; i < maxRetries; i++ {
		// Preparar Payload
		payload := map[string]string{"message": prompt}
		jsonPayload, _ := json.Marshal(payload)

		// Cliente con Timeout para no quedarse colgado
		client := &http.Client{Timeout: 45 * time.Second}
		resp, err := client.Post(targetURL, "application/json", bytes.NewBuffer(jsonPayload))

		if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
			log.Printf("⚠️ Intento %d fallido. Reintentando...", i+1)
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(2 * time.Second)
			continue
		}

		var result struct {
			Reply string `json:"reply"`
		}

		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if err == nil && strings.TrimSpace(result.Reply) != "" {
			return result.Reply
		}

		log.Printf("⚠️ Respuesta vacía en intento %d", i+1)
		time.Sleep(1 * time.Second)
	}

	return "Lo siento, mi cerebro está saturado. ¿Podrías intentar en un momento?"
}