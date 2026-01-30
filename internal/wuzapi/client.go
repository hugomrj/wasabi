package wuzapi

import (
    "bytes"
    "encoding/json"
    "log"
    "fmt"
    "net/http"
    "os"
    "wasabi/internal/models"
)

func SendMessage(token string, phone string, body string) error {
    url := os.Getenv("WUZAPI_URL") + "/chat/send/text"
    
    // Asegúrate de que este JSON sea exactamente lo que Wuzapi espera
    msg := map[string]interface{}{
        "Phone": phone,
        "Body":  body,
    }
    
    jsonData, _ := json.Marshal(msg)
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Token", token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("❌ Error de red al enviar a Wuzapi: %v", err)
        return err
    }
    defer resp.Body.Close()

    // ESTO ES LO MÁS IMPORTANTE AHORA:
    log.Printf("📡 Respuesta de Wuzapi al enviar: Código %d", resp.StatusCode)
    return nil
}    