# Skill: pubsub-push-subscriber-blueprint

## Objetivo
Definir el esqueleto base de un subscriber HTTP Push en Go.

## Checklist minimo
- Endpoint `POST /` para Pub/Sub push.
- Parseo envelope y decode base64 de `message.data`.
- Parse de evento de negocio.
- Logs estructurados (`messageId`, `eventType`, `orderId`).
- ACK rules:
  - `200` en ignorados, duplicados y exito.
  - `500` en fallo critico para retry.

## Arquitectura recomendada
- `domain`, `application`, `interfaces`, `infrastructure`, `platform`.

## Entregables minimos
- health check
- handler push
- caso de uso principal
- workflow deploy Cloud Run
