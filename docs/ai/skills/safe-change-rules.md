# Skill: safe-change-rules

## Checklist antes de merge
- ¿Se mantiene semantica de ACK 200/500?
- ¿Se mantiene idempotencia por `eventType+orderId`?
- ¿Se cubren duplicados sin reproceso?
- ¿Se agregaron tests de regresion?
- ¿Se actualizaron docs operativas?

## Regla critica
Ningun cambio debe hacer que Telegram determine exito/fallo del procesamiento principal.
