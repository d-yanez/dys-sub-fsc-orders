# Skill: telegram-operational-notifications

## Objetivo
Estandarizar mensajes Telegram de operacion para subscribers.

## Reglas
- Mensaje final solo al terminar procesamiento.
- Nunca bloquear flujo principal si Telegram falla.
- No enviar Telegram en duplicados.

## Contenido sugerido
- `eventType`, `orderId`, `orderNumber`.
- `resultado` (`SUCCESS|PARTIAL_SUCCESS|FAILED`).
- `phase` y `error/warning` si aplica.
- Datos copiables con `<code>`.
- Link operativo adicional (si aplica), por ejemplo stock por SKU.
