# Skill: idempotency-event-logs-v1

## Objetivo
Implementar idempotencia simple y trazabilidad por evento usando `event_logs`.

## Patron v1
- 1 documento por `idempotencyKey`.
- Key sugerida: `eventType + orderId`.

## Ciclo de vida
- `RECEIVED` -> `PROCESSING` -> `SUCCESS|PARTIAL_SUCCESS|FAILED`.
- `processed=true` solo al cierre exitoso/partial.

## Regla duplicados
- Si `processed=true` en la key: `200`, sin reproceso, sin notificacion.

## Campos minimos
- `_id`, `idempotencyKey`, `messageId`, `eventType`, `orderId`, `payloadHash`, `status`, `attempts`, `phase`, `processed`, `processedAt`, `updatedAt`.
