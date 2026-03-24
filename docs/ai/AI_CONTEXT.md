# AI_CONTEXT

## Proyecto
- Servicio: `dys-sub-fsc-orders`
- Stack: Go 1.22, Cloud Run, Pub/Sub push, MongoDB.
- Rol: subscriber/orquestador de persistencia para eventos FSC `onOrderCreated`.

## Objetivo funcional v1
- Procesar `onOrderCreated`.
- Persistir `orders`, `order_items`, `event_logs`.
- Soportar retries e idempotencia simplificada.
- Notificar resultado final a Telegram.

## No objetivos v1
- No exponer CRUD de ordenes.
- No procesar otros eventos FSC.
- No implementar lease lock complejo.
