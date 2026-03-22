# Event Flow (PR4)

1. Recibe `POST /` de Pub/Sub push.
2. Valida envelope (`message.data`).
3. Decodifica base64.
4. Parsea evento Falabella.
5. Normaliza `eventType` (`eventType || event`).
6. Si no es `onOrderCreated`, responde `200` e ignora.
7. Si es `onOrderCreated` y trae `payload.OrderId`, consulta `/order/{orderId}`.
8. Consulta `/orderItems/{orderId}`.
9. Enriquecimiento por SKU (`/sku/{sku}`) para thumbnail (no critico).
10. Registra/actualiza `event_logs` con clave `falabella:onOrderCreated:{orderId}`.
11. Si `processed=true`, marca duplicado e ignora (`200`).
12. Persiste `orders` y `order_items` en MongoDB.
13. Marca `event_logs` como `SUCCESS` o `PARTIAL_SUCCESS` con `processed=true`.
14. Si falla una dependencia critica (FSC order/items o Mongo), marca `FAILED`, `processed=false` y responde `500` (retry habilitado).
