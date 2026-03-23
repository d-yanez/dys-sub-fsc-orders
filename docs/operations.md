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
- Formato operativo al finalizar:
  - `SUCCESS`
  - `PARTIAL_SUCCESS`
  - `FAILED`
- Si Telegram falla, no revierte persistencia principal.
- En duplicados (`duplicate_ignored`) no se envia Telegram.

## OIDC
- Middleware stub disponible para activar validacion en PRs siguientes:
  - `OIDC_VALIDATION_ENABLED`
  - `OIDC_AUDIENCE`
  - `OIDC_ALLOWED_EMAIL`
