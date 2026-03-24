# SDD (Software Design Description)

## Proposito del servicio
Persistir ordenes Falabella al recibir `onOrderCreated` desde Pub/Sub push, con idempotencia, trazabilidad y notificacion operacional.

## Responsabilidades
- Ingesta de envelope Pub/Sub push.
- Validacion basica del evento y filtro por tipo.
- Enriquecimiento de orden/items desde FSC.
- Persistencia en MongoDB (`orders`, `order_items`).
- Gestion de `event_logs` para idempotencia y retries.
- Telegram final (`SUCCESS`, `PARTIAL_SUCCESS`, `FAILED`).

## Entradas
- HTTP `POST /` (Pub/Sub push envelope).
- `message.data` base64 con evento Falabella.
- Formato soportado para `payload.OrderId`: string y numerico.

## Salidas
- HTTP response al push (`200` o `500` segun regla).
- Documentos Mongo actualizados.
- Logs estructurados para observabilidad.
- Mensaje Telegram final (opcional por switch).

## Dependencias
- FSC API (`/order/{orderId}`, `/orderItems/{orderId}`, `/sku/{sku}`)
- MongoDB
- Telegram API
- Google Cloud Run + Pub/Sub

## Notas de diseno
- v1 usa idempotencia simplificada sin lease lock complejo.
- Duplicado procesado: `200`, sin reproceso y sin Telegram.
- Fallo de thumbnail o Telegram: no tumba flujo principal (puede ser `PARTIAL_SUCCESS`).
