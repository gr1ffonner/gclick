package main

import (
	"gclick/internal/app"
	"gclick/internal/db"
	"gclick/pkg/config"
	"gclick/pkg/logging"
)

//	@title			API for test task
//	@version		1.0
//	@description	API for test task
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8888
//	@BasePath	/api/v1

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
