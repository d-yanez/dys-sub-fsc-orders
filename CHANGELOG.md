# Changelog

## [0.1.0] - 2026-03-22
- Scaffold inicial de `dys-sub-fsc-orders`.
- Estructura clean por capas (`domain`, `application`, `interfaces`, `infrastructure`, `platform`).
- Endpoint `POST /` para Pub/Sub push con parseo base64 y filtro `onOrderCreated`.
- Workflow de deploy a Cloud Run (`cr-dys-sub-fsc-orders-prd`).
