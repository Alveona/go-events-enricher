package system

import (
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

type Container struct {
	Clickhouse clickhouse.ChConn
}

func New(clickhouseDB clickhouse.ChConn) *Container {
	return &Container{
		Clickhouse: clickhouseDB,
	}
}
