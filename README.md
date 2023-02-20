## go-events-enricher

This is a simple (yet production-ready) microservice for async proxy analytic events to clickhouse

### Quickstart
`make docker-server-run`, runs clickhouse container with go application, migrations applied

### API Docs
Check `swagger-doc/` in repository  
Payload is expected to be in a format of escaped json list, separated by `\n` newline character  
```{\"client_time\":\"2000-01-01 00:00:00\",\"device_id\":\"00000000-0000-0000-0000-000000000001\",\"device_os\":\"iOS\",\"session\":\"s1\",\"sequence\":1,\"event\":\"example_1\",\"param_int\":123,\"param_str\":\"string1\"}\n{\"client_time\":\"2000-01-02 00:00:00\",\"device_id\":\"00000000-0000-0000-0000-000000000002\",\"device_os\":\"Android\",\"session\":\"s2\",\"sequence\":2,\"event\":\"example_2\",\"param_int\":1234,\"param_str\":\"string2\"}```

### Structure
Service is a typical DDD application
- `handlers`, HTTP endpoint handlers
- `processors`, Business-logics layer
- `storages`, Data repositories
- `generated`, go-swagger generated stuff
- `metrics`, prometheus metrics container

This service uses the following [Base service template](https://github.com/Alveona/go-base-service) as basement
