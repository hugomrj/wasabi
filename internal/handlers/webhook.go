package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Responder 200 OK inmediatamente para evitar timeouts en Wuzapi
	w.WriteHeader(http.StatusOK)

	r.ParseForm()
	rawJSON := r.FormValue("jsonData")
	if rawJSON == "" {
		return
	}

	// 2. Procesar de forma asíncrona (Goroutine)
	go func(data string) {
		var payload WebhookPayload
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			log.Printf("❌ Error decodificando payload: %v", err)
			return
		}

		if payload.EventData.Info.IsFromMe {
			return
		}

		// Extraer texto (robusto para diferentes tipos de mensajes)
		prompt := payload.EventData.Message.Conversation
		if prompt == "" {
			prompt = payload.EventData.Message.ExtendedText.Text
		}

		if prompt == "" {
			return
		}

		// Limpiar y formatear el número del remitente
		remitente := payload.EventData.Info.Sender
		if payload.EventData.Info.SenderAlt != "" {
			remitente = payload.EventData.Info.SenderAlt
		}
		remitente = strings.Split(strings.Split(remitente, "@")[0], ":")[0]

		log.Printf("📩 Consultando IA para [%s]...", remitente)

		// 3. Llamada a la IA (Puede tardar, pero no bloquea el servidor)
		respuestaIA := getAIResponse(prompt)

		// 4. Enviar respuesta final a WhatsApp
		token := "USER_TOKEN_1"
		err := wuzapi.SendMessage(token, remitente, respuestaIA)
		if err != nil {
			log.Printf("❌ Error enviando a %s: %v", remitente, err)
		} else {
			log.Printf("✅ Respuesta enviada con éxito a %s", remitente)
		}
	}(rawJSON)
}

// getAIResponse conecta con el servidor Python para obtener la respuesta de la IA
func getAIResponse(prompt string) string {
	url := "https://japo.click/charlette/ask"

	payload := map[string]string{"message": prompt}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("❌ Error al conectar con Python: %v", err)
		return "Lo siento, tuve un problema al conectar con mi cerebro artificial."
	}
	defer resp.Body.Close()

	var result struct {
		Reply string `json:"reply"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("❌ Error decodificando respuesta de la IA: %v", err)
		return "Recibí una respuesta extraña de la IA."
	}

	return result.Reply
}