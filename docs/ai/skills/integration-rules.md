# Skill: integration-rules

## FSC
- `GET /order/{orderId}` y `GET /orderItems/{orderId}` son criticos.
- `GET /sku/{sku}` es no critico (si falla: warning/partial).

## Mongo
- Colecciones: `orders`, `order_items`, `event_logs`.
- Cuidar indices e impacto en consultas futuras por `orderNumber` y `orderId`.

## Pub/Sub push
- `200` evita retry.
- `500` habilita retry para fallo critico.

## Telegram
- Mensaje operativo final.
- Si falla envio, registrar error pero no revertir persistencia principal.
