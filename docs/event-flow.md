# Event Flow (PR2)

1. Recibe `POST /` de Pub/Sub push.
2. Valida envelope (`message.data`).
3. Decodifica base64.
4. Parsea evento Falabella.
5. Normaliza `eventType` (`eventType || event`).
6. Si no es `onOrderCreated`, responde `200` e ignora.
7. Si es `onOrderCreated` y trae `payload.OrderId`, responde `200` con `accepted`.

Nota: en PR2 no hay persistencia ni notificacion final; solo ingestion y decision de enrutamiento.
