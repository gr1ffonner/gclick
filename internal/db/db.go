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
	log  logging.Logger
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
		log:  logger,
	}, nil
}

func (c *ClickhouseWriter) Insert(ctx context.Context, a Article) error {
	err := c.conn.Exec(ctx, "INSERT INTO events (eventID,eventType, userID, eventTime, payload) VALUES (rand64(), ?, ?, ?, ?)", a.EventType, a.UserID, a.EventTime, a.Payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClickhouseWriter) GetbyEventTypeANDTime(ctx context.Context, evType string, evTime string) (Article, error) {
	var a Article
	r := c.conn.QueryRow(ctx, "SELECT eventID, eventType, userID, toString(eventTime), payload FROM events WHERE eventType = ? AND eventTime = ?", evType, evTime)
	// convert eventtime to string

	err := r.Scan(&a.EventID, &a.EventType, &a.UserID, &a.EventTime, &a.Payload)
	if err != nil {
		c.log.Info(err)
	}
	return a, nil
}
