# dys-sub-fsc-orders

Subscriber HTTP Push (Pub/Sub) para eventos de Falabella Seller Center.

En esta etapa (PR2) el servicio implementa:
- bootstrap Go y estructura por capas
- endpoint push `POST /`
- parseo de envelope Pub/Sub (`message.data` base64)
- normalizacion de evento (`eventType || event`)
- filtro de `onOrderCreated`
- punto preparado para validacion OIDC (`middleware` stub)

## Endpoints
- `GET /health`
- `POST /`

## Variables de entorno base
- `PORT` (default `8080`)
- `SERVICE_NAME` (default `dys-sub-fsc-orders`)
- `ENVIRONMENT` (default `local`)
- `LOG_LEVEL` (default `INFO`)
- `PUBSUB_SUBSCRIPTION_NAME` (opcional)
- `OIDC_VALIDATION_ENABLED` (default `false`)
- `OIDC_AUDIENCE` (opcional)
- `OIDC_ALLOWED_EMAIL` (opcional)

## Ejecutar local
```bash
go run ./cmd/api
```

## Tests
```bash
go test ./...
```

## Scope v1 (pendiente en PRs siguientes)
- Integracion FSC (`/order`, `/orderItems`, `/sku`)
- Mongo (`orders`, `order_items`, `event_logs`)
- idempotencia y retries
- Telegram final operativo
