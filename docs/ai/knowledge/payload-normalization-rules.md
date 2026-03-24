# Skill: payload-normalization-rules

## Objetivo
Evitar fallos por variaciones de contrato en eventos externos.

## Reglas practicas
- Soportar alias de nombre (`eventType` y `event`).
- Soportar tipos mixtos donde tenga sentido (ej: `OrderId` string o numerico).
- Normalizar entrada a formato interno estable (trim + string).

## Resultado esperado
- Menos `invalid_event_json` en produccion.
- Menos dependencias fragiles del productor.
