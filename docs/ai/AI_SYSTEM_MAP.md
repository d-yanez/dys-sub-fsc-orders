# AI_SYSTEM_MAP

## Entrada
- Pub/Sub push HTTP -> `POST /` (Cloud Run subscriber).

## Core
- Handler Pub/Sub -> UseCase `ProcessOnOrderCreated`.
- Validacion/filtro de evento.
- Idempotencia y ciclo de vida en `event_logs`.

## Integraciones
- FSC API: `/order`, `/orderItems`, `/sku`.
- MongoDB: `orders`, `order_items`, `event_logs`.
- Telegram Bot API: mensaje final operativo.

## Salida operacional
- ACK HTTP al push.
- Persistencia de orden/items/log.
- Logging estructurado para observabilidad.
