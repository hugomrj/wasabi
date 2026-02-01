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
		respuestaIA := GetExternalResponse(prompt)

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


// GetExternalResponse envía un prompt a un servicio externo y devuelve la respuesta procesada.
func GetExternalResponse(prompt string) string {
	// Podrías incluso pasar la URL como parámetro si quieres que sea 100% genérica
	const targetURL = "https://japo.click/charlette/ask"

	// 1. Empaquetar el mensaje
	payload := map[string]string{"message": prompt}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Error al serializar JSON: %v", err)
		return "Error interno: no se pudo procesar el formato del mensaje."
	}

	// 2. Realizar la petición con un tiempo límite (opcional pero recomendado)
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("❌ Error de red/conexión: %v", err)
		return "No se pudo establecer conexión con el servicio externo."
	}
	defer resp.Body.Close()

	// 3. Verificar que el servidor destino respondió correctamente
	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ El servicio externo devolvió código: %d", resp.StatusCode)
		return "El servicio externo encontró un error al procesar la solicitud."
	}

	// 4. Decodificar la respuesta esperada
	var result struct {
		Reply string `json:"reply"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("❌ Error al decodificar la respuesta: %v", err)
		return "La respuesta recibida no tiene un formato válido."
	}

	// 5. Validar que no llegue vacío
	finalText := strings.TrimSpace(result.Reply)
	if finalText == "" {
		log.Printf("⚠️ El servicio devolvió una respuesta vacía")
		return "No se obtuvo una respuesta válida del servicio."
	}

	return finalText
}