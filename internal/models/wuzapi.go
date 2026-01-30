package models

type WuzapiRequest struct {
    Event string `json:"event"`
    Data  struct {
        From   string `json:"from"`
        Body   string `json:"body"`
        FromMe bool   `json:"fromMe"` // <--- ESTO EVITA EL BUCLE
    } `json:"data"`
}

type TextPayload struct {
    Phone string `json:"Phone"`
    Body  string `json:"Body"`
}