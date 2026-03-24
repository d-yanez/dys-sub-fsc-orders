# ARCHITECTURE

## Proposito tecnico
`dys-sub-fsc-orders` es un subscriber HTTP Push de Pub/Sub que procesa eventos Falabella `onOrderCreated`, enriquece datos con APIs FSC, persiste en MongoDB y envia notificacion final a Telegram.

## Componentes principales
- `cmd/api`: bootstrap, wiring de dependencias y servidor HTTP.
- `internal/interfaces/http`: router, middleware y handler Pub/Sub.
- `internal/application`: DTOs y caso de uso `ProcessOnOrderCreatedUseCase`.
- `internal/domain`: entidades, puertos y servicios de dominio (idempotency key).
- `internal/infrastructure`: clientes FSC/Telegram y repositorios Mongo.
- `internal/platform`: config, logger y utilitarios.

## Flujo general
1. Pub/Sub push envia `POST /` con envelope.
2. Handler decodifica `message.data` (base64) y parsea evento.
3. Use case filtra solo `onOrderCreated`.
4. Se controla idempotencia en `event_logs`.
5. Se consulta FSC (`/order`, `/orderItems`, `/sku`).
6. Se persisten `orders` y `order_items`.
7. Se cierra `event_logs` en `SUCCESS`, `PARTIAL_SUCCESS` o `FAILED`.
8. Se envia Telegram final (no bloqueante).

## Dependencias relevantes
- MongoDB (`falabellaDB`): `orders`, `order_items`, `event_logs`.
- FSC API (`dys-api-fsc`): datos de orden, items y thumbnail.
- Telegram Bot API: notificacion operativa final.
- Google Cloud Pub/Sub Push + Cloud Run.

## Integraciones externas confirmadas
- Entrada: Pub/Sub push HTTP.
- Salida: FSC API, MongoDB, Telegram.
- Deploy: GitHub Actions + Workload Identity + Cloud Run.
