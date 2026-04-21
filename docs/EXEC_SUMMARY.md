# Executive Summary — dys-sub-fsc-orders

## Qué es
`dys-sub-fsc-orders` es un subscriber HTTP Push (Pub/Sub) que procesa eventos Falabella (`onOrderCreated`), enriquece datos contra FSC y persiste `orders`/`order_items` en Mongo.

## Patrón principal
- **Subscriber event-driven** con idempotencia y política ACK/retry.

## Qué hace / qué NO hace
- Qué hace:
  - Parsea envelope Pub/Sub.
  - Filtra eventos de interés.
  - Consulta APIs FSC (`order`, `orderItems`, `sku`).
  - Persiste órdenes/items y `event_logs`.
  - Envía Telegram final (no bloqueante).
- Qué no hace:
  - No emite DTE ni orquesta boletas.
  - No expone API de negocio para frontend.

## Endpoints
- `GET /health`
- `POST /` (push Pub/Sub)

## Reglas críticas
- Idempotencia obligatoria por evento/orden.
- ACK `200` para duplicados y casos procesados.
- `500` cuando corresponde retry por falla crítica.
- Persistir pricing de items (`price`, `paidPrice`, `shippingAmount`) para flujos downstream DTE.

## Operación
- Revisar primero `event_logs` y luego `orders` / `order_items`.
- Correlacionar por `messageId`, `eventType`, `orderId`.
- Telegram es señal operativa, no fuente de verdad.

## Checklist de evolución segura
- [ ] Mantener política ACK/retry explícita por tipo de error.
- [ ] Mantener idempotencia y consistencia de `event_logs`.
- [ ] Mantener contrato de persistencia requerido por consumidores downstream.
- [ ] Actualizar `event-flow`, `operations` y esta ficha al cambiar flujo.
