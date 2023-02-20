#!/bin/sh
set -ex

go install github.com/pressly/goose/v3/cmd/goose@latest

cd "${MIGRATIONS_DIR}" || (code=$? && echo "Migrations directory '${WORK_DIR}' does not exist or has been not provided!" && exit ${code})
goose clickhouse "${EVENTS_ENRICHER_CLICKHOUSE_MIGRATION_DSN}" up || (code=$? && echo "Migrations cannot be applied!" && exit ${code})

cd "${WORK_DIR}" || (code=$? && echo "Working directory '${WORK_DIR}' does not exist or has been not provided!" && exit ${code})
go run -ldflags "-s -w" main.go