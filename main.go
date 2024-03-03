package main

import (
	"cobaApp/config"
	"cobaApp/database"
	"cobaApp/logger"
	"cobaApp/server"
	tracing "cobaApp/tracing"
	"fmt"
	"github.com/opentracing/opentracing-go"
)

func main() {
	fmt.Println("haloo semua")

	// load config
	cfg := config.NewConfig()

	fmt.Println(cfg.GetConfig().Database.Name)

	// define log console
	log := logger.NewConsoleLog()

	// define tracing
	tracer, closer := tracing.GenerateTracing(cfg, log, "cobaApp")
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	// connect to database
	db := database.ConnectDatabase(cfg, log)
	defer db.Close()

	appServer := server.NewAppServer(db, cfg)
	if err := appServer.RunServer(); err != nil {
		log.Fatalf("cant start app : %v", err)
	} else {
		log.Info("success run app!")
	}
}
