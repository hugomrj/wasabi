# 🌿 Wasabi – Wuzapi Webhook Processor


**Wasabi** es una capa de orquestación y middleware de alto rendimiento desarrollada en Go. Su función principal es actuar como puente inteligente entre **Wuzapi** (que gestiona la conexión con WhatsApp) y servicios externos de **IA**.

Wasabi no se conecta directamente a WhatsApp; en su lugar, recibe los eventos de las múltiples instancias de Wuzapi, los procesa y los distribuye a sus respectivos motores de IA de forma eficiente y aislada.


### 🔄 Flujo de comunicación:
`WhatsApp 📱 <-> Wuzapi 🔌 <-> Wasabi (Go) 🌿 <-> Tu IA 🤖`


## 🚀 Características

- **Multi-instancia dinámico:** Gestiona múltiples clientes mediante rutas variables (`/webhook/{id_instancia}`).
- **Procesamiento asíncrono:** Utiliza *goroutines* para procesar mensajes en segundo plano sin bloquear el flujo del webhook.
- **Control de concurrencia:** Implementa un sistema de semáforos para gestionar la carga de peticiones hacia la IA externa.
- **Identificación en logs:** Cada evento está etiquetado con el nombre de la instancia para una depuración rápida.
- **Agnóstico:** Permite añadir nuevos clientes editando únicamente el archivo `.env`, sin necesidad de recompilar el binario.

## 📂 Estructura del proyecto

```text
.
├── cmd/main.go               # Punto de entrada del servidor
├── internal/
│   ├── handlers/             # Manejadores de rutas y lógica de IA
│   ├── models/               # Estructuras de datos (payloads del webhook)
│   └── wuzapi/               # Cliente para envío de mensajes
├── .env                      # Variables de entorno (tokens y URLs)
└── go.mod                    # Dependencias del proyecto
```

## ⚙️ Configuración (.env)

Wasabi utiliza un sistema de mapeo dinámico basado en el ID de la instancia recibido en la URL.
Por cada cliente, debés definir su token de Wuzapi y la URL de su servicio de IA correspondiente.

### Puerto donde corre Wasabi
```
WASABI_PORT=3000
```


### URL de Wazapi
```
WUZAPI_URL=http://localhost:8080 
```






```

# --- CONFIGURACIÓN DE CLIENTES ---
# Formato:
# ID_INSTANCIA=TOKEN_WUZAPI
# ID_INSTANCIA_URL=URL_IA

# Ejemplo: cliente "ventas"

ventas=TU_TOKEN_WUZAPI_AQUI
ventas_URL=https://tu-ia.com/ventas/ask

# Ejemplo: cliente "soporte"
soporte=OTRO_TOKEN_WUZAPI
soporte_URL=https://tu-ia.com/soporte/ask
```

## 📡 Uso del webhook
Para que los mensajes lleguen a Wasabi, debés configurar la URL del webhook en cada instancia de Wuzapi usando el ID de instancia definido en el .env.

URL del webhook:
```
http://TU_IP_O_DOMINIO:3000/webhook/{ID_INSTANCIA}
Registro vía cURL
curl -X POST http://localhost:8080/instance/set \
  -H "token: TOKEN_DE_LA_INSTANCIA" \
  -H "Content-Type: application/json" \
  -d '{
    "webhook_url": "http://TU_IP:3000/webhook/ventas"
  }'
```

### 🛠️ Despliegue
Compilación manual
```
go build -o wasabi cmd/main.go
./wasabi
```

Nota técnica
El sistema vincula automáticamente la ruta /webhook/{id} con las variables {id} y {id}_URL definidas en el entorno.
Si el ID no existe en el archivo .env o la configuración está incompleta, Wasabi rechazará la petición y lo notificará en los logs.


---
