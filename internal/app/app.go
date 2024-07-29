package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gr1ffonner/fintest/internal/config"
	"github.com/gr1ffonner/fintest/pkg/db"
	"github.com/gr1ffonner/fintest/pkg/logging"
	"github.com/julienschmidt/httprouter"

	httpSwagger "github.com/swaggo/http-swagger"
)

const layout = "2006-01-02 15:04:05"

// @Summary		Insert material
// @Description	Insert material
// @Tags			materials
// @Accept			json
// @Produce		text/html
// @Param data body db.Material true "The input material struct"
// @Success		200
// @Failure		400	{object}	error	"Bad request"
// @Failure		404	{object}	error	"Not found"
// @Failure		500	{object}	error	"Internal Server Error"
// @Router			/material [post]
func insertMaterial(db *db.DB, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		validate := validator.New()

		var m db.Material

		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = validate.Struct(m)
		if err != nil {
			logger.Error(err)

			validationErrors, ok := err.(validator.ValidationErrors)
			if ok {
				for _, fieldError := range validationErrors {
					http.Error(w, fmt.Sprintf("You should provide %s", fieldError.Field()), http.StatusBadRequest)
				}
				return
			}
			return
		}

		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()

		err = db.CreateMaterial(m)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Material created successfully")
	}
}

// @Summary		Get material by UUID
// @Description	Get material by UUID
// @Tags			materials
// @Param			id	path	string	true	"material UUID"
// @Success		200	{object}	db.Material
// @Failure		400	{object}	error	"Bad request"
// @Failure		404	{object}	error	"Not found"
// @Failure		500	{object}	error	"Internal Server Error"
// @Router			/material/{id} [get]
func getMaterialByUUID(db *db.DB, logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		id := httprouter.ParamsFromContext(r.Context()).ByName("id")
		if id == "" {
			http.Error(w, "Bad request, UUID not provided", http.StatusBadRequest)
			return
		}

		m, err := db.GetMaterialByUUID(id)
		if err != nil {
			logger.Error(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	}
}

// @Summary		healthcheck
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

func StartApp(cfg config.Config, logger logging.Logger, db *db.DB) {
	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/material", insertMaterial(db, logger))
	router.HandlerFunc("GET", "/api/v1/material/:id", getMaterialByUUID(db, logger))
	router.HandlerFunc("GET", "/api/v1/healthcheck", healthcheck())
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	logger.Infof("server is working on %s%s", cfg.Server.Host, cfg.Server.Port)
	err := http.ListenAndServe(cfg.Server.Port, router)
	if err != nil {
		logger.Fatal(err)
	}
}
