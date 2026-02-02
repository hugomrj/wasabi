    package handlers

    import "net/http"

    // MapRoutes organiza todas las APIs del sistema.
    // MapRoutes registra únicamente las rutas necesarias para la operación.
    func MapRoutes(mux *http.ServeMux) {
        
        // Ruta de diagnóstico (indispensable para saber si el proceso corre)
        mux.HandleFunc("/api/health/ping", PingHandler)

        // Ruta única de integración con Wuzapi
        //  mux.HandleFunc("/api/v1/wuzapi/webhook", WebhookHandler)

        // Usamos el método POST porque Wuzapi envía los mensajes así.
        mux.HandleFunc("POST /webhook/{instancia}", WebhookHandler)


    }