# Changelog

## [0.1.0] - 2026-03-22
- Scaffold inicial de `dys-sub-fsc-orders`.
- Estructura clean por capas (`domain`, `application`, `interfaces`, `infrastructure`, `platform`).
- Endpoint `POST /` para Pub/Sub push con parseo base64 y filtro `onOrderCreated`.
- Workflow de deploy a Cloud Run (`cr-dys-sub-fsc-orders-prd`).
- Integracion FSC para consultar `/order`, `/orderItems`, `/sku`.
- Persistencia core en MongoDB para `orders` y `order_items`.
- `event_logs` con ciclo de vida `RECEIVED -> PROCESSING -> SUCCESS|PARTIAL_SUCCESS|FAILED`.
- Idempotencia simplificada v1 por `eventType + orderId`.
- Duplicados procesados se ignoran (`200`) sin reproceso.
