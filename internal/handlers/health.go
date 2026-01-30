package handlers

import "net/http"

// PingHandler responde con un simple "pong" para verificar que el servidor está vivo.
// Debe empezar con Mayúscula para que sea visible si se necesitara fuera del paquete,
// aunque aquí al estar en el mismo paquete 'handlers' basta con que exista.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}


func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}