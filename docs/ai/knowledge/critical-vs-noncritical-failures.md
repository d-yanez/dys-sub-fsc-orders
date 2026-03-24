# Skill: critical-vs-noncritical-failures

## Objetivo
Separar fallos criticos de no criticos para controlar retry correctamente.

## Criticos (deben devolver 500)
- Fetch de datos esenciales del proveedor.
- Persistencia principal en base de datos.

## No criticos (no deben tumbar proceso)
- Enriquecimiento opcional (ej: thumbnail).
- Notificacion externa (ej: Telegram).

## Resultado esperado
- `SUCCESS`: core completo sin warnings.
- `PARTIAL_SUCCESS`: core OK con warning no critico.
- `FAILED`: fallo critico, `processed=false`, retry habilitado.
