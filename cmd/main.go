package main

import (
	"gclick/internal/app"
	"gclick/internal/db"
	"gclick/pkg/config"
	"gclick/pkg/logging"
)

func main() {
	// Setup logging
	logger := logging.GetLogger("trace")
	logger.Info("logger is working")
	// Read config
	cfg := config.GetConfig()
	logger.Info("config is working")

	// Setup db
	db, err := db.NewClickhouseWriter(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("db is working")

	logger.Info("Ci introduced")
	// Start app
	app.StartApp(*cfg, logger, db)

}
