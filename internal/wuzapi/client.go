package wuzapi

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "wasabi/internal/models"
)

func SendMessage(token string, phone string, text string) error {
    baseURL := os.Getenv("WUZAPI_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080"
    }

    apiURL := fmt.Sprintf("%s/chat/send/text", baseURL)

    payload := models.TextPayload{
        Phone: phone,
        Body:  text,
    }

    jsonData, _ := json.Marshal(payload)

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    // Pasamos el Token recibido para que Wuzapi sepa qué instancia usar
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Token", token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}