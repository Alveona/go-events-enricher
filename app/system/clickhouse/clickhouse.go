package clickhouse

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	_ "github.com/ClickHouse/clickhouse-go" //nolint:golint
	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

func New(config *Config) (ChConn, error) {
	dsn := config.DSN + "&max_execution_time=" + strconv.Itoa(config.MaxExecutionTime)

	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db := &chDb{
		master:             conn,
		healthCheckTimeout: config.HealthCheckTimeout,
	}

	db.master.SetMaxIdleConns(config.MaxIdleConns)
	db.master.SetMaxOpenConns(config.MaxOpenConns)
	db.master.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}

type chDb struct {
	master             *sql.DB
	healthCheckTimeout time.Duration
}

func (db *chDb) Master() *sql.DB {
	return db.master
}

func (db *chDb) PingMaster() bool {
	ctx, cancel := context.WithTimeout(context.Background(), db.healthCheckTimeout)
	defer cancel()

	if err := db.master.PingContext(ctx); err != nil {
		logrus.Errorf("ping Clickhouse failed: %v", err)
		return false
	}
	return true
}
