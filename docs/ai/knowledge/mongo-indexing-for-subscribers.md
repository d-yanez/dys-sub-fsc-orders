# Skill: mongo-indexing-for-subscribers

## Objetivo
Definir indices minimos para consultas y trazabilidad.

## Recomendacion base
- `orders`: unique `orderNumber`, indice `status+createdAt`.
- `order_items`: indice `orderId`, indice compuesto `orderId+sku`, indice `status`.
- `event_logs`: indice `messageId`, `orderId+updatedAt`, `status+updatedAt`.

## Regla
Alinear indices a casos de uso reales de lectura y debug operativo, no solo a escritura.
