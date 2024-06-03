package app

import (
	"context"
	"encoding/json"
	"fmt"
	_ "gclick/docs"
	"gclick/internal/db"
	"gclick/pkg/config"
	"gclick/pkg/logging"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter" // Using httprouter for better route handling

	httpSwagger "github.com/swaggo/http-swagger"
)

const layout = "2006-01-02 15:04:05"

// @Summary		Insert article
//
// @Description	Insert article
// @Tags			events
// @Accept			json
// @Success		200
// @Failure		400	{object}	error	"Bad request"
// @Failure		404	{object}	error	"Not found"
// @Failure		500	{object}	error	"Internal Server Error"
// @Router			/event [post]
func insertArticle(clickhouse db.ClickhouseWriter, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var a db.Article
		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		t, err := time.Parse(layout, a.EventTime)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		strt := t.Format(time.RFC3339Nano)
		strt = strings.ReplaceAll(strt, "T", " ")
		strt = strings.ReplaceAll(strt, "Z", "")

		err = clickhouse.Insert(ctx, a)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "everything is fine %s", strt)
	}
}

// @Summary		Get articles by type and time
//
// @Description	Get articles by type and time
// @Tags			events
// @Param			eventType	path		string	true	"event type"
// @Param			eventTime	path		string	true	"event time"
// @Success		200			{object}	db.Article
// @Failure		400			{object}	error	"Bad request"
// @Failure		404			{object}	error	"Not found"
// @Failure		500			{object}	error	"Internal Server Error"
// @Router			/event/{eventType}/{eventTime} [get]
func getArticlesByTypeaANDTime(clickhouse db.ClickhouseWriter, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		params := httprouter.ParamsFromContext(ctx)
		evType := params.ByName("eventType")
		evTime := params.ByName("eventTime")

		a, err := clickhouse.GetbyEventTypeANDTime(ctx, evType, evTime)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(a)
	}
}

// @Summary		healthcheck
//
// @Description	healthcheck endpoint
// @Tags			healthcheck
// @Success		200
// @Failure		400	{object}	error	"Bad request"
// @Failure		404	{object}	error	"Not found"
// @Failure		500	{object}	error	"Internal Server Error"
// @Router			/healthcheck [get]
func healthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// send 200 status code
		w.WriteHeader(http.StatusOK)
	}
}

func StartApp(cfg config.Config, logger logging.Logger, db *db.ClickhouseWriter) {
	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/event", insertArticle(*db, logger))
	router.HandlerFunc("GET", "/api/v1/event/:eventType/:eventTime", getArticlesByTypeaANDTime(*db, logger))
	router.HandlerFunc("GET", "/api/v1/healthcheck", healthcheck())
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Infof("server is working on %s%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, router)
	if err != nil {
		logger.Fatal(err)
	}
}
