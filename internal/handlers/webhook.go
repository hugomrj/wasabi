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

	err := json.Unmarshal([]byte(rawJSON), &wrapper)

	if err != nil {
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
        token := "USER_TOKEN_1" 
        
        log.Printf("📩 Consultando IA para %s...", remitente)

        // CAMBIO: En lugar de "Capturé...", llamamos a la función puente
        mensajeDeRespuesta := llamarAlServidorIA(msg.Conversation)

		// ESTO ES LO QUE DEBES AGREGAR:
		log.Printf("🤖 IA respondió: [%s]", mensajeDeRespuesta)

        err := wuzapi.SendMessage(token, remitente, mensajeDeRespuesta)

		if err != nil {
			log.Printf("❌ Error enviando: %v", err)
		} else {
			log.Println("✅ Respuesta enviada con éxito")
		}
	}

	w.WriteHeader(http.StatusOK)
}


// Esta función va fuera de WebhookHandler, al final del archivo
func llamarAlServidorIA(textoUsuario string) string {
    // La dirección de tu servidor Python
    url := "https://japo.click/tu-bot/ask"

    // 1. Creamos el sobre (JSON) con el parámetro "message" que espera Python
    payload := map[string]string{"message": textoUsuario}
    jsonPayload, _ := json.Marshal(payload)

    // 2. Enviamos el sobre a la otra computadora
    resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
    if err != nil {
        log.Printf("❌ Error al contactar servidor IA: %v", err)
        return "Lo siento, tengo problemas para conectarme con mi cerebro."
    }
    defer resp.Body.Close()

    // 3. Abrimos la respuesta de Python para sacar el "reply"
    var result struct {
        Reply string `json:"reply"`
    }
    json.NewDecoder(resp.Body).Decode(&result)

    return result.Reply
}