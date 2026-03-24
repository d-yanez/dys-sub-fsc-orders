# SKILL_INDEX

## Skills bootstrap
- `skills/project-ai-context-bootstrap.md`
  - Bootstrapping de contexto IA y docs base.

## Skills de dominio actual
- `skills/service-overview.md`
  - Resumen rapido del servicio y limites de alcance.
- `skills/event-processing.md`
  - Flujo detallado de procesamiento y estados.
- `skills/safe-change-rules.md`
  - Reglas para cambios seguros sin romper operacion.
- `skills/integration-rules.md`
  - Reglas practicas para FSC, Mongo, Telegram y Pub/Sub.

## Knowledge reusable para futuros subscribers
- `knowledge/pubsub-push-subscriber-blueprint.md`
  - Plantilla base de subscriber push en Go.
- `knowledge/idempotency-event-logs-v1.md`
  - Patron de idempotencia simple con `event_logs`.
- `knowledge/critical-vs-noncritical-failures.md`
  - Matriz de decision para `SUCCESS|PARTIAL_SUCCESS|FAILED`.
- `knowledge/telegram-operational-notifications.md`
  - Estandar de mensajes Telegram operativos.
- `knowledge/payload-normalization-rules.md`
  - Reglas para tolerar contratos variables sin romper parseo.
- `knowledge/cloud-run-pubsub-oidc-hardening.md`
  - Checklist IAM/OIDC/WIF para Cloud Run + Pub/Sub push.
- `knowledge/mongo-indexing-for-subscribers.md`
  - Indices recomendados para persistencia/event logs.
- `knowledge/production-smoke-and-replay-runbook.md`
  - Smoke y replay operativo en produccion.
