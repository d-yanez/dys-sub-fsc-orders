# Changelog

## [0.1.3] - 2026-04-01
- Ajuste quirúrgico en Telegram `onOrderCreated`: ahora se envía **1 mensaje por ítem persistido** (N mensajes para N ítems), manteniendo la misma estructura textual operativa.
- Cada mensaje usa los datos del ítem correspondiente (`orderItemId`, `sku`, `item`) y su link de stock por SKU.
- Se agregó soporte de `ThumbnailURL` por ítem para que cada mensaje adjunte su imagen del SKU (si existe).
- Tests actualizados para validar el comportamiento multi-mensaje y contenido por ítem.
- Cobertura adicional agregada para validar explícitamente escenarios de 1, 2 y N ítems (`N=5`) con envío por ítem.

## [0.1.2] - 2026-04-01
- Fix quirúrgico en notificación Telegram de `onOrderCreated`: ahora incluye el detalle de **todos** los ítems persistidos (orderItemId, sku, qty, nombre), no solo el primer ítem.
- Se mantiene compatibilidad con el formato existente (`itemsPersistidos`, resumen y primer ítem) para no romper consumidores actuales.
- Nuevo test: `TestProcessSuccessTelegramIncludesAllPersistedItems` para validar el caso de orden con 3 ítems.

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
