package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"wasabi/internal/wuzapi"
)

// WebhookPayload se mantiene igual (puedes moverlo a internal/models/wuzapi.go si prefieres)
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

var iaSemaphore = make(chan struct{}, 1)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer token dinámico de la URL (Go 1.22+)
	instancia := r.PathValue("instancia")
	if instancia == "" {
		log.Printf("⚠️ Webhook recibido sin instancia en la URL")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// LOG VISUAL PARA DEPURACIÓN
    log.Println("========================================")
    log.Printf(" ID INSTANCIA DETECTADA: %s ", instancia)
    log.Println("========================================")	

	// 2. Responder 200 OK de inmediato a Wuzapi
	w.WriteHeader(http.StatusOK)

	// 3. Procesar Formulario
	r.ParseForm()
	rawJSON := r.FormValue("jsonData")
	if rawJSON == "" {
		return
	}


	token := os.Getenv(instancia)    
    targetURL := os.Getenv(instancia + "_URL") 
    // Validación de seguridad para ambos
    if token == "" || targetURL == "" {
        log.Printf("❌ ERROR: Configuración incompleta para '%s' en .env (Falta token o URL)", instancia)
        return 
    }



	// 4. Lanzar proceso en segundo plano pasando el token de instancia
	go func(data string,  instancia string, token string, targetURL string) {
		var payload WebhookPayload
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
            // USAR 'instancia' aquí en lugar de 'token'
            log.Printf("❌ [%s] Error decodificando payload: %v", instancia, err)
            return
        }


		// No responderse a uno mismo
		if payload.EventData.Info.IsFromMe {
			return
		}

		// Extraer el texto del mensaje
		prompt := payload.EventData.Message.Conversation
		if prompt == "" {
			prompt = payload.EventData.Message.ExtendedText.Text
		}
		if prompt == "" {
			return
		}

		// Limpiar número de teléfono
		remitente := payload.EventData.Info.Sender
		if payload.EventData.Info.SenderAlt != "" {
			remitente = payload.EventData.Info.SenderAlt
		}
		remitente = strings.Split(strings.Split(remitente, "@")[0], ":")[0]



		// Obtener respuesta de la IA (Gestionado por semáforo)
		respuestaIA := GetExternalResponse(prompt, targetURL)

		// 5. USAR EL TOKEN EXTRAÍDO para responder
		err := wuzapi.SendMessage(token, remitente, respuestaIA)
		if err != nil {
            // USAR 'instancia' aquí para el log
            log.Printf("❌ [%s] Error enviando a %s: %v", instancia, remitente, err)
        } else {
            // USAR 'instancia' aquí para el log
            log.Printf("✅ [%s] Respuesta enviada a %s", instancia, remitente)
        }
	}(rawJSON, instancia, token, targetURL)
}



func GetExternalResponse(prompt string, targetURL string) string {
	log.Printf("⏳ Mensaje en espera de turno para IA...")
	iaSemaphore <- struct{}{}
	defer func() { <-iaSemaphore }()

	// const targetURL = "https://japo.click/charlette/ask"
	const maxRetries = 2

	for i := 0; i < maxRetries; i++ {
		payload := map[string]string{"message": prompt}
		jsonPayload, _ := json.Marshal(payload)

		client := &http.Client{Timeout: 45 * time.Second}
		resp, err := client.Post(targetURL, "application/json", bytes.NewBuffer(jsonPayload))

		if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
			log.Printf("⚠️ Reintento %d IA fallido", i+1)
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
	}

	return "Lo siento ¿Podrías intentar en un momento?"
}