package app

import (
	"context"
	"encoding/json"
	"fmt"
	"gclick/internal/db"
	"gclick/pkg/config"
	"gclick/pkg/logging"
	"net/http"
	"time"
)

const layout = "2006-01-02 15:04:05"

func getArticlesById(clickhouse *db.ClickhouseWriter, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		// unmarshal request
		var a db.Article

		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			logger.Error(err)
		}
		// Parse the value of the a.EventTime field as a time.Time value
		t, err := time.Parse("2006-01-02 15:04:05", a.EventTime)
		if err != nil {
			logger.Error(err)
		}
		a.EventTime = t.Format(time.RFC3339Nano)

		err = clickhouse.Insert(ctx, a)
		if err != nil {
			logger.Error(err)
		}
		fmt.Fprintf(w, "everything is fine")
	}
}

func welcome(logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "welcome")
		logger.Info("hit welcome")
	}
}

func StartApp(cfg config.Config, logger logging.Logger, db *db.ClickhouseWriter) {
	// Setup server
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/event", getArticlesById(db, logger))
	// mux.HandleFunc("GET /trending", getTrending(db, logger))
	mux.HandleFunc("GET /", welcome(logger))
	logger.Infof("server is working on %s%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, mux)
	if err != nil {
		logger.Fatal(err)
	}
}
