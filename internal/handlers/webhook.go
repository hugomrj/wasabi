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
            // Captura también mensajes con formato o respuestas
            ExtendedText struct {
                Text string `json:"text"`
            } `json:"extendedTextMessage"`
        } `json:"Message"`
        Type string `json:"type"`
    } `json:"event"`
}

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Respondemos 200 OK inmediatamente para liberar la conexión de entrada
    w.WriteHeader(http.StatusOK)

    r.ParseForm()
    rawJSON := r.FormValue("jsonData")
    if rawJSON == "" {
        return
    }

    // 2. Procesamos en segundo plano (Asíncrono)
    go func(data string) {
        var payload WebhookPayload
        if err := json.Unmarshal([]byte(data), &payload); err != nil {
            log.Printf("❌ Error decodificando payload: %v", err)
            return
        }

        // Ignorar si el mensaje lo enviamos nosotros
        if payload.EventData.Info.IsFromMe {
            return
        }

        // Extraer el texto (priorizar conversación, luego texto extendido)
        texto := payload.EventData.Message.Conversation
        if texto == "" {
            texto = payload.EventData.Message.ExtendedText.Text
        }

        if texto == "" {
            log.Println("⚠️ Mensaje sin texto recibido")
            return
        }

        // Limpiar remitente
        remitente := payload.EventData.Info.Sender
        if payload.EventData.Info.SenderAlt != "" {
            remitente = payload.EventData.Info.SenderAlt
        }
        remitente = strings.Split(strings.Split(remitente, "@")[0], ":")[0]

        log.Printf("📩 Procesando mensaje de [%s]: %s", remitente, texto)

        // 3. Llamada a la IA (Aquí Go esperará pacientemente el tiempo que haga falta)
        respuestaIA := llamarAlServidorIA(texto)

        // 4. Enviar de vuelta a WhatsApp
        token := "USER_TOKEN_1" 
        err := wuzapi.SendMessage(token, remitente, respuestaIA)
        if err != nil {
            log.Printf("❌ Error enviando respuesta a %s: %v", remitente, err)
        } else {
            log.Printf("✅ Respuesta enviada con éxito a %s", remitente)
        }
    }(rawJSON)
}