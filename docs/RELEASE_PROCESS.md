# RELEASE_PROCESS

## Flujo actual
1. Cambios en branch.
2. Merge/push a `main`.
3. GitHub Actions (`.github/workflows/deploy.yml`) ejecuta:
   - `go mod tidy`
   - `go build ./...`
   - `go test ./...`
   - `go vet ./...`
   - build/push imagen
   - deploy a Cloud Run `cr-dys-sub-fsc-orders-prd`

## Requisitos de despliegue
- WIF (Workload Identity) habilitado para repo.
- SA deploy con permisos de deploy.
- SA runtime con permisos de ejecucion necesarios.
- Servicio Cloud Run privado (`--no-allow-unauthenticated`).

## Estrategia prudente
- Validar run en GitHub Actions y revision `Ready` en Cloud Run.
- Ejecutar smoke test con `onOrderCreated` real.
- Confirmar:
  - logs de subscriber
  - persistencia en Mongo (`orders`, `order_items`, `event_logs`)
  - Telegram (si switch habilitado)

## Rollback
- Reasignar trafico a revision previa de Cloud Run.
- Mantener trazabilidad con `messageId` y `orderId` en logs/event logs.
