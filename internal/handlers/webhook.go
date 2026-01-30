package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"wasabi/internal/models"
	"wasabi/internal/wuzapi"
)

// WebhookHandler procesa las peticiones HTTP entrantes desde Wuzapi.
// Implementa la interfaz http.HandlerFunc estándar de Go.
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Restricción de Método: Solo permitimos POST según el estándar de webhooks.
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// 2. Gestión de Memoria: Cerramos el cuerpo de la petición al finalizar la función.
	defer r.Body.Close()

	// 3. Preparación del Modelo: Usamos el molde definido en internal/models.
	var payload models.WuzapiRequest

	// 4. Decodificación de Datos: Transformamos el JSON entrante a un Struct de Go.
	// Se usa json.NewDecoder por ser más eficiente con flujos de datos (streams).
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Error decodificando webhook: %v", err)
		http.Error(w, "JSON malformado", http.StatusBadRequest)
		return
	}

	// 5. Filtrado de Eventos: Wuzapi envía varios tipos de eventos;
	// solo procesamos "Message" para evitar respuestas en bucle o innecesarias.
	if payload.Event == "Message" {
		remitente := payload.Data.From
		mensaje := payload.Data.Body

		// Validamos que el mensaje no esté vacío antes de procesar.
		if mensaje != "" {
			log.Printf("Procesando mensaje de: %s", remitente)

			// 6. Respuesta Saliente: Utilizamos el cliente parametrizado.
			// Se ignora el error con '_' por simplicidad, aunque en producción debería loguearse.
			_ = wuzapi.SendMessage(remitente, "Recibido en Wasabi (Go Puro) ✅")
		}
	}

	// 7. Confirmación de Recepción: Respondemos con HTTP 200 OK.
	// Esto es crucial para que Wuzapi no considere la entrega como fallida y reintente.
	w.WriteHeader(http.StatusOK)
}


