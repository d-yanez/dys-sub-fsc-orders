# Operations (PR2)

## Health check
- `GET /health`

## Pub/Sub push
- Endpoint: `POST /`
- ACK en PR2: siempre `200` para envelopes parseables, incluyendo eventos ignorados.

## OIDC
- Middleware stub disponible para activar validacion en PRs siguientes:
  - `OIDC_VALIDATION_ENABLED`
  - `OIDC_AUDIENCE`
  - `OIDC_ALLOWED_EMAIL`
