# PROJECT_MEMORY

## Decisiones importantes
- Idempotencia v1: un doc por `idempotencyKey` (`eventType + orderId`).
- Duplicado procesado: ACK `200`, sin reproceso y sin Telegram.
- Fallo critico: ACK `500` para retry de Pub/Sub.
- `orders._id = orderId`, `order_items._id = orderItemId`.
- `orders.orderNumber` con indice unico para busqueda futura.

## Hechos operativos recientes
- Se ajusto parser de evento para aceptar `OrderId` numerico.
- Se agrego link `Ver stock bodega` en mensaje Telegram final.
