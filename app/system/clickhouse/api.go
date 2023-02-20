package clickhouse

import "database/sql"

type ChConn interface {
	Master() *sql.DB
	PingMaster() bool
}
