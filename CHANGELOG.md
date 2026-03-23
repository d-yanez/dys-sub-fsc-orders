# Changelog

## [0.1.1] - 2026-03-23
- Telegram subscriber ahora agrega link operativo por SKU al final del mensaje: `Ver stock bodega`.
- Nuevo env `STOCK_VIEW_BASE_URL` con default `https://dy-api-utils-785293986978.us-central1.run.app/stock/view`.
- Link construido de forma segura con `PathEscape` del SKU y omision cuando no hay SKU.
- Tests unitarios agregados para inclusion/omision/escape del link.

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
- Telegram final del subscriber con formato operativo para `SUCCESS|PARTIAL_SUCCESS|FAILED`.
- `SWITCH_FSC_ORDER_TELEGRAM` para habilitar/deshabilitar notificaciones sin afectar procesamiento principal.
