package db

import (
	"context"
	"fmt"
	"gclick/pkg/config"
	"gclick/pkg/logging"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// DB is a struct to hold the database connection.
type ClickhouseWriter struct {
	conn driver.Conn
}

func NewClickhouseWriter(logger logging.Logger, cfg *config.Config) (*ClickhouseWriter, error) {
	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database.Database,
			Username: cfg.Database.Username,
			Password: cfg.Database.Password,
		},
		Debug: true,
	})
	if err != nil {
		logger.Error(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		logger.Error(err)
	}

	return &ClickhouseWriter{
		conn: conn,
	}, nil
}

func (c *ClickhouseWriter) Insert(ctx context.Context, a Article) error {
	err := c.conn.Exec(ctx, "INSERT INTO events (eventType, userID, eventTime, payload) VALUES (?, ?, ?, ?)", a.EventType, a.UserID, a.EventTime, a.Payload)
	if err != nil {
		return err
	}
	return nil
}
