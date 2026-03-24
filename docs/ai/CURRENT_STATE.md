# CURRENT_STATE

## Estado funcional actual
- Subscriber en produccion en Cloud Run.
- Procesa `onOrderCreated` end-to-end.
- `OrderId` soporta string y numerico en payload.
- Telegram incluye link operativo "Ver stock bodega" por SKU.

## Riesgos abiertos conocidos
- Dependencia fuerte de disponibilidad FSC API.
- OIDC middleware esta en modo stub (hardening pendiente).
- Infra aun no declarada totalmente como IaC (Terraform pendiente).

## Prioridades sugeridas
1. Consolidar runbook de operacion/incidentes.
2. Migrar infra a Terraform con impersonation.
3. Endurecer seguridad OIDC en push.
