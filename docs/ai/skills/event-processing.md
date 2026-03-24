# Skill: event-processing

## Uso
Para analizar o depurar el flujo de un evento `onOrderCreated`.

## Secuencia
1. Parse envelope Pub/Sub.
2. Decodifica `message.data`.
3. Evalua evento y `orderId`.
4. `event_logs`: `RECEIVED -> PROCESSING`.
5. Consulta FSC (`/order`, `/orderItems`, `/sku`).
6. Persiste `orders` y `order_items`.
7. Cierra `event_logs` en `SUCCESS|PARTIAL_SUCCESS|FAILED`.
8. Envia Telegram final si corresponde.

## Claves de debug
- `messageId`
- `orderId`
- `eventType`
- `phase`
