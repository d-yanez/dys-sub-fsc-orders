# Skill: production-smoke-and-replay-runbook

## Objetivo
Validar end-to-end en produccion y reprocesar casos puntuales de forma controlada.

## Smoke minimo
1. Disparar evento de prueba real.
2. Confirmar publicacion del productor (topic/messageId).
3. Confirmar recepcion y ACK del subscriber.
4. Confirmar persistencia Mongo.
5. Confirmar notificacion final.

## Replay puntual
- Verificar estado en `event_logs`.
- Si aplica, limpiar key/idempotency del caso.
- Reinyectar evento.
- Validar nuevo `messageId` y cierre en `SUCCESS|PARTIAL_SUCCESS`.
