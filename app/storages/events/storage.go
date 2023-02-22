package events

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Workiva/go-datastructures/queue"
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"

	"github.com/pkg/errors"

	"github.com/Alveona/go-events-enricher/app/entities"
	"github.com/Alveona/go-events-enricher/app/system/clickhouse"
)

const eventsTable = "events.logs"

type StorageMetrics interface {
	QueryDuration(duration time.Duration, labels ...string)
}

type BufferMetrics interface {
	BufferWriteIOInc()
	BufferReadIOInc()
}

type CHStorage struct {
	conn           clickhouse.ChConn
	storageMetrics StorageMetrics
	bufferMetrics  BufferMetrics
	rl             ratelimit.Limiter
	q              *queue.RingBuffer
	bufferCfg      clickhouse.BufferConfig
}

func NewCHRepo(conn clickhouse.ChConn, storageMetrics StorageMetrics, bufferMetrics BufferMetrics, bufferCfg clickhouse.BufferConfig) *CHStorage {
	return &CHStorage{
		conn:           conn,
		storageMetrics: storageMetrics,
		bufferMetrics:  bufferMetrics,
		rl:             ratelimit.New(bufferCfg.Ratelimit),
		q:              queue.NewRingBuffer(bufferCfg.BufferSize),
		bufferCfg:      bufferCfg,
	}
}

func (r *CHStorage) ProcessInsertEvents(ctx context.Context, events []*entities.EventDTO) error {
	if !r.bufferCfg.WithBuffer {
		return r.insertEvents(ctx, events)
	}
	for {
		added, _ := r.q.Offer(events)
		if added {
			r.bufferMetrics.BufferWriteIOInc()
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (r *CHStorage) StartQueueProducing() {
	go func() {
		for {
			ctx := context.Background()
			item, err := r.q.Poll(1)
			if err != nil {
				time.Sleep(r.bufferCfg.DequeueDelay)
				continue
			}
			r.bufferMetrics.BufferReadIOInc()

			r.rl.Take()
			events, ok := item.([]*entities.EventDTO)
			if !ok {
				logrus.Errorf("An invalid event type: %+v", events)
				continue
			}

			err = r.insertEvents(ctx, events)
			if err != nil {
				logrus.Errorf("Failed to insert events from buffer: %+v", err)
				continue
			}
		}
	}()

}

func (r *CHStorage) StopQueueProducing() error {
	done := make(chan bool, 1)
	logrus.Info("Clickhouse buffer has been interrupted")
	go func() {
		for r.q.Len() != 0 {
			time.Sleep(r.bufferCfg.BufferWaitRetries)
		}
		logrus.Info("Clickhouse buffer is empty")
		done <- true
	}()

	select {
	case <-done:
		logrus.Info("Clickhouse buffer has been gracefully stopped")
		return nil
	case <-time.After(r.bufferCfg.DequeueTimeout):
		return errors.New("Reached the timeout while shutting down the Clickhouse queue")
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

func (r *CHStorage) insertEvents(ctx context.Context, events []*entities.EventDTO) error {

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
		r.storageMetrics.QueryDuration(time.Since(queryStart), "InsertEvents")
	}

	return nil
}
