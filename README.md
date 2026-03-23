# dys-sub-fsc-orders

Subscriber HTTP Push (Pub/Sub) para eventos de Falabella Seller Center.

En esta etapa (PR4) el servicio implementa:
- bootstrap Go y estructura por capas
- endpoint push `POST /`
- parseo de envelope Pub/Sub (`message.data` base64)
- normalizacion de evento (`eventType || event`)
- filtro de `onOrderCreated`
- cliente FSC (`/order`, `/orderItems`, `/sku`)
- persistencia core en MongoDB (`orders`, `order_items`)
- `event_logs` con idempotencia simplificada (`processed`, `status`, `attempts`)
- duplicados por `eventType + orderId` ignorados con ACK `200`
- retries habilitados para fallos criticos (`500`)
- enriquecimiento de thumbnail por SKU (degradado a warning si falla)
- notificacion Telegram final (`SUCCESS`, `PARTIAL_SUCCESS`, `FAILED`) con switch
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
- `MONGODB_URI` (requerida para correr el servicio)
- `MONGODB_DB_NAME` (default `falabellaDB`)
- `DYS_API_FSC_BASE_URL` (default `https://dys-api-fsc-785293986978.us-central1.run.app`)
- `DYS_API_FSC_API_KEY` (recomendada)
- `HTTP_TIMEOUT_MS_FSC` (default `5000`)
- `SWITCH_FSC_ORDER_TELEGRAM` (`true|false`, default `false`)
- `TELEGRAM_BOT_TOKEN`
- `TELEGRAM_CHAT_ID`
- `TELEGRAM_TIMEOUT_MS` (default `5000`)

## Ejecutar local
```bash
go run ./cmd/api
```

## Tests
```bash
go test ./...
```

## Scope v1 pendiente (PRs siguientes)
- cierre deploy productivo

