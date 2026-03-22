# Operations (PR2)

## Health check
- `GET /health`

## Pub/Sub push
- Endpoint: `POST /`
- ACK en PR4:
  - `200` para duplicados procesados o casos `SUCCESS/PARTIAL_SUCCESS`
  - `500` para fallo critico (habilita retry)

## OIDC
- Middleware stub disponible para activar validacion en PRs siguientes:
  - `OIDC_VALIDATION_ENABLED`
  - `OIDC_AUDIENCE`
  - `OIDC_ALLOWED_EMAIL`
