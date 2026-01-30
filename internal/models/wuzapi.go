package models

// TextPayload es lo que enviamos HACIA Wuzapi
type TextPayload struct {
	Phone string `json:"Phone"`
	Body  string `json:"Body"`
}

// WuzapiRequest es lo que recibimos DESDE Wuzapi vía Webhook
type WuzapiRequest struct {
	Event string `json:"event"`
	Data  struct {
		From string `json:"from"`
		Body string `json:"body"`
	} `json:"data"`
}