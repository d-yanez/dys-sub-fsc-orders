# Operations (PR2)

## Health check
- `GET /health`

## Pub/Sub push
- Endpoint: `POST /`
- ACK en PR4:
  - `200` para duplicados procesados o casos `SUCCESS/PARTIAL_SUCCESS`
  - `500` para fallo critico (habilita retry)

## Telegram
- Variable de control: `SWITCH_FSC_ORDER_TELEGRAM`
- Base del link operativo por SKU: `STOCK_VIEW_BASE_URL` (default `https://dy-api-utils-785293986978.us-central1.run.app/stock/view`)
- Formato operativo al finalizar:
  - `SUCCESS`
  - `PARTIAL_SUCCESS`
  - `FAILED`
- Si hay SKU del primer item, se agrega al final: `ver stock bodega: Ver stock bodega` (link clickeable).
- Si Telegram falla, no revierte persistencia principal.
- En duplicados (`duplicate_ignored`) no se envia Telegram.

## OIDC
- Middleware stub disponible para activar validacion en PRs siguientes:
  - `OIDC_VALIDATION_ENABLED`
  - `OIDC_AUDIENCE`
  - `OIDC_ALLOWED_EMAIL`

## Backfill historico de `orders.financial` y `orders.addresses`
- Dry run:
  - `./scripts/backfill_orders_financial.sh --dry-run=true --limit=100 --batch-size=50`
- Ejecucion real:
  - `./scripts/backfill_orders_financial.sh --dry-run=false --limit=1000 --batch-size=100 --from-order-id=1140000000`
