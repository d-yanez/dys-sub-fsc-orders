# ENGINEERING_PRACTICES

## Cambios seguros
- Mantener responsabilidad acotada: este servicio procesa solo `onOrderCreated`.
- No romper semantica de ACK:
  - `200` en ignorados/duplicados/exito.
  - `500` solo en fallo critico para permitir retry.
- Preservar idempotencia v1 por `eventType + orderId` en `event_logs`.
- Telegram nunca debe romper persistencia principal.

## Convenciones tecnicas observadas
- Go por capas: `domain`, `application`, `interfaces`, `infrastructure`, `platform`.
- Logs estructurados (`slog`) con claves operativas (`messageId`, `orderId`, `eventType`, `phase`).
- Config por env vars con defaults prudentes en `internal/platform/config`.
- Repositorios con contratos en `domain/ports` e implementacion en `infrastructure`.

## Buenas practicas para contribuir
- Agregar/actualizar tests junto con cambios de comportamiento.
- Evitar cambios de schema sin documentar impacto en indices y consultas.
- Documentar cambios operativos en `docs/CHANGELOG.md` y runbooks existentes.
- Mantener mensajes de error accionables para troubleshooting en Cloud Logging.
