# Architecture

Servicio Go orientado a Clean Architecture:

- `domain`: entidades y contratos de dominio
- `application`: casos de uso y DTOs
- `interfaces`: capa HTTP (ingest Pub/Sub push)
- `infrastructure`: adaptadores concretos (pendientes en PR3+)
- `platform`: config, logger y utilitarios transversales

En PR2 el foco es ingestion y enrutamiento del evento. Persistencia e integraciones externas quedan para PR3+.
