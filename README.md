# 🌿 Wasabi - Wuzapi Webhook Processor

**Wasabi** es un middleware de alto rendimiento desarrollado en Go, diseñado para actuar como puente entre **Wuzapi** (WhatsApp API) y servicios de Inteligencia Artificial externos. Su arquitectura está optimizada para entornos **multi-instancia**, permitiendo gestionar múltiples cuentas de WhatsApp con configuraciones de IA independientes desde un único servidor.

## 🚀 Características

- **Multi-instancia Dinámico:** Gestiona múltiples clientes mediante rutas variables (`/webhook/{id_instancia}`).
- **Procesamiento Asíncrono:** Utiliza Goroutines para procesar mensajes en segundo plano sin bloquear el flujo del webhook.
- **Control de Concurrencia:** Implementa un sistema de semáforos para gestionar la carga de peticiones hacia la IA externa.
- **Identificación en Logs:** Cada evento está etiquetado con el nombre de la instancia para una depuración rápida.
- **Agnóstico:** Añade nuevos clientes simplemente editando el archivo `.env`, sin necesidad de recompilar el binario.

## 📂 Estructura del Proyecto

```text
.
├── cmd/main.go               # Punto de entrada del servidor
├── internal/
│   ├── handlers/             # Manejadores de rutas y lógica de IA
│   ├── models/               # Estructuras de datos (Webhook payloads)
│   └── wuzapi/               # Cliente para envío de mensajes
├── .env                      # Variables de entorno (Tokens y URLs)
└── go.mod                    # Dependencias del proyecto
```

⚙️ Configuración (.env)
Wasabi utiliza un sistema de mapeo dinámico basado en el ID de la instancia que llega por la URL. Por cada cliente, debes añadir su Token de Wuzapi y su URL de IA correspondiente en el archivo .env:

Fragmento de código
# Puerto donde corre Wasabi
WASABI_PORT=3000

# --- CONFIGURACIÓN DE CLIENTES ---
# Formato: 
# NOMBRE_ID=TOKEN_WUZAPI
# NOMBRE_ID_URL=URL_IA_CORRESPONDIENTE

# Ejemplo para un cliente llamado 'ventas'
ventas=TU_TOKEN_WUZAPI_AQUI
ventas_URL=[https://tu-ia.com/ventas/ask](https://tu-ia.com/ventas/ask)

# Ejemplo para un cliente llamado 'soporte'
soporte=OTRO_TOKEN_WUZAPI
soporte_URL=[https://tu-ia.com/soporte/ask](https://tu-ia.com/soporte/ask)
📡 Uso del Webhook
Para que los mensajes lleguen a Wasabi, debes configurar la URL del webhook en cada instancia de Wuzapi utilizando el ID definido en tu archivo de configuración:

URL del Webhook: http://TU_IP_O_DOMINIO:3000/webhook/{ID_INSTANCIA}

Registro vía CURL:
```Bash
curl -X POST http://localhost:8080/instance/set \
  -H "token: TOKEN_DE_LA_INSTANCIA" \
  -H "Content-Type: application/json" \
  -d '{
    "webhook_url": "http://TU_IP:3000/webhook/ventas"
  }'
```

🛠️ Despliegue
Compilación Manual
Si deseas compilar el binario en tu entorno:

```Bash
go build -o wasabi cmd/main.go
./wasabi
```

Nota Técnica: El sistema vincula automáticamente la ruta /webhook/xyz con las variables xyz y xyz_URL definidas en el entorno. Si el ID no existe en el .env o la configuración está incompleta, Wasabi rechazará la petición y lo notificará en los logs.