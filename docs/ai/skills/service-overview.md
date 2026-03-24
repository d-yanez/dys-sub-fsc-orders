# Skill: service-overview

## Uso
Cuando se necesite explicar rapido que hace `dys-sub-fsc-orders`.

## Resumen operativo
- Entrada: Pub/Sub push (`POST /`).
- Procesa: solo `onOrderCreated`.
- Enriquecimiento: FSC (`order`, `orderItems`, `sku`).
- Persistencia: Mongo (`orders`, `order_items`, `event_logs`).
- Notificacion: Telegram final (opcional por switch).

## Limites
- No CRUD publico.
- No otros eventos FSC en v1.
