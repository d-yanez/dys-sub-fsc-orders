# Skill: cloud-run-pubsub-oidc-hardening

## Objetivo
Checklist de seguridad/autenticacion para push Pub/Sub a Cloud Run.

## Checklist
- Cloud Run privado (`--no-allow-unauthenticated`).
- Push subscription con OIDC y audience correcto.
- `roles/run.invoker` para SA de push.
- `roles/iam.serviceAccountTokenCreator` para servicio Pub/Sub sobre SA de push.
- WIF de GitHub con `attributeCondition` que incluya repo correcto.
- Deploy SA con `iam.serviceaccounts.actAs` sobre runtime SA.

## Errores comunes
- `unauthorized_client` por `attributeCondition` mal configurada.
- `actAs denied` por falta de `roles/iam.serviceAccountUser`.
