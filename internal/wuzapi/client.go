package wuzapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"wasabi/internal/models"
)


// SendMessage ahora recibe el apiToken como tercer argumento
func SendMessage(phone, body, apiToken string) error {
    baseURL := os.Getenv("WUZAPI_URL")

    payload := models.TextPayload{
        Phone: phone,
        Body:  body,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    apiURL := fmt.Sprintf("%s/chat/send/text", baseURL)
    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    // ¡AQUÍ ESTÁ EL CAMBIO! Usamos el token que llega por parámetro
    req.Header.Set("Token", apiToken) 
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}