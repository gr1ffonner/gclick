package app

import (
	"context"
	"encoding/json"
	"fmt"
	"gclick/internal/db"
	"gclick/pkg/config"
	"gclick/pkg/logging"
	"net/http"
	"strings"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

const layout = "2006-01-02 15:04:05"

func insertArticle(clickhouse db.ClickhouseWriter, logger logging.Logger) http.HandlerFunc {
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

		strt := t.Format(time.RFC3339Nano)
		strt = strings.ReplaceAll(strt, "T", " ")
		strt = strings.ReplaceAll(strt, "Z", "")
		logger.Info(a.EventTime)

		err = clickhouse.Insert(ctx, a)
		if err != nil {
			logger.Error(err)
		}
		fmt.Fprintf(w, "everything is fine %s", a.EventTime)
	}
}

func getArticlesByTypeaANDTime(clickhouse db.ClickhouseWriter, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)

		defer cancel()
		evType := r.PathValue("eventType")
		evTime := r.PathValue("eventTime")
		a, err := clickhouse.GetbyEventTypeANDTime(ctx, evType, evTime)
		if err != nil {
			logger.Error(err)
		}
		fmt.Fprintf(w, "%+v", a)
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
	mux.HandleFunc("POST /api/event", insertArticle(*db, logger))
	mux.HandleFunc("GET /api/event/{eventType}/{eventTime}", getArticlesByTypeaANDTime(*db, logger))
	mux.HandleFunc("GET /", welcome(logger))
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8888/swagger/doc.json")))
	logger.Infof("server is working on %s%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, mux)
	if err != nil {
		logger.Fatal(err)
	}
}
