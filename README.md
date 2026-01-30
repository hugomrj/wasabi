# 🌿 Wasabi - Wuzapi Webhook Processor

**Wasabi** es un middleware ligero desarrollado en Go diseñado para recibir eventos de **Wuzapi** y procesarlos de manera eficiente. Está optimizado para entornos multi-instancia, permitiendo gestionar múltiples números de WhatsApp desde un solo punto de enlace.

## 🚀 Características

- **Multi-instancia:** Detecta automáticamente qué instancia envía el mensaje mediante headers de Token.
- **Arquitectura Limpia:** Separación clara entre modelos, manejadores y lógica de cliente.
- **Instalación Automatizada:** Incluye un script en Python para despliegue rápido en servidores Linux.
- **Servicio de Sistema:** Configurado para correr como un servicio de `systemd` (24/7).

## 📂 Estructura del Proyecto

```text
.
├── cmd/wasabi/main.go        # Punto de entrada del servidor
├── internal/
│   ├── handlers/             # Lógica de las rutas (Webhook, Health)
│   ├── models/               # Definición de estructuras JSON
│   └── wuzapi/               # Cliente para enviar mensajes a Wuzapi
├── .env                      # Configuración de entorno (no incluido en git)
├── go.mod                    # Dependencias de Go
└── wasabi_installer.py       # Script de instalación automática
```


🛠️ Instalación en Servidor (Ubuntu)
1. Requisitos Previos
Go 1.21+ instalado (sudo apt install golang-go)

Python 3

2. Despliegue Rápido
Utiliza el instalador incluido para desplegar en /srv/wasabi:

```Bash
python3 wasabi_installer.py
```
El script se encargará de:

Clonar el repositorio.

Crear el archivo .env.

Compilar el binario de Go.

Crear y activar el servicio en systemd.

⚙️ Configuración (.env)
El archivo .env debe contener las siguientes variables:

WUZAPI_URL: Dirección base donde corre tu API de Wuzapi (ej. http://localhost:8080).

WASABI_PORT: Puerto donde escuchará este webhook (ej. 3000).

📡 Uso del Webhook
Para que Wuzapi envíe mensajes a Wasabi, debes registrar el webhook en cada instancia:

```Bash

curl -X POST http://localhost:8080/webhook \
  -H "Token: TU_USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "webhook": "http://TU_IP_SERVIDOR:3000/webhook",
    "events": ["Message"]
  }'
```  
📊 Monitoreo y Logs
Para ver la actividad del webhook en tiempo real:

```Bash
journalctl -u wasabi -f
```