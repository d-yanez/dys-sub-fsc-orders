# WORKING_RULES

- No inventar rutas ni servicios fuera del monorepo.
- Mantener separacion por capas (`domain/application/interfaces/infrastructure/platform`).
- Antes de cambios, validar impacto en:
  - ACK/retry
  - idempotencia
  - trazabilidad (`messageId`, `orderId`, `phase`)
- No acoplar este servicio con CRUD futuro.
- Telegram nunca debe bloquear persistencia principal.
- Mantener docs operativas actualizadas tras cambios productivos.
