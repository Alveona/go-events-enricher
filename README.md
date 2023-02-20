## go-events-enricher

This is a simple (yet production-ready) microservice for async proxy analytic events to clickhouse

### Quickstart
`make docker-server-run`, runs clickhouse container with go application, migrations applied

### API Docs
Check `swagger-doc/` in repository

### Structure
Service is a typical DDD application
- `handlers`, HTTP endpoint handlers
- `processors`, Business-logics layer
- `storages`, Data repositories
- `generated`, go-swagger generated stuff
- `metrics`, prometheus metrics container

This service uses the following [Base service template](https://github.com/Alveona/go-base-service) as basement