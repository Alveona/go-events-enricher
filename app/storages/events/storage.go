package events

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/pkg/errors"

	"github.com/Alveona/go-events-enricher/app/entities"
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

const eventsTable = "events.logs"

type Metrics interface {
	QueryDuration(duration time.Duration, labels ...string)
}

type CHStorage struct {
	conn    clickhouse.ChConn
	metrics Metrics
}

func NewCHRepo(conn clickhouse.ChConn, m Metrics) *CHStorage {
	return &CHStorage{
		conn:    conn,
		metrics: m,
	}
}

func (r *CHStorage) GetTransaction(ctx context.Context) (*sql.Tx, error) {
	return r.conn.Master().BeginTx(ctx, nil)
}

func prepareInsertEventQuery(event *entities.EventDTO) (string, []interface{}, error) {
	sb := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	rawRequest := sb.Insert(eventsTable)
	st, err := clickhouse.NewStommer(event)
	if err != nil {
		return "", []interface{}{}, errors.Wrap(err, "create stommer")
	}
	rawRequest = rawRequest.Columns(st.Columns...).Values(st.Values...)

	return rawRequest.ToSql()
}

func (r *CHStorage) InsertEvents(ctx context.Context, events []*entities.EventDTO) error {

	for _, event := range events {
		// Clickhouse driver requires "transaction" to present to be able to insert values.
		// However, clickhouse doesn't support transaction, so we do fake transaction each time
		queryStart := time.Now()
		tx, err := r.GetTransaction(ctx)
		if err != nil {
			return errors.WithStack(err)
		}
		query, args, err := prepareInsertEventQuery(event)
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			return errors.WithStack(err)
		}
		err = tx.Commit()
		if err != nil {
			return errors.WithStack(err)
		}
		r.metrics.QueryDuration(time.Since(queryStart), "InsertEvents")
	}

	return nil
}
