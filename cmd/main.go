package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/Chandra5468/cfp-Products-Service/cmd/httpapi"
	"github.com/Chandra5468/cfp-Products-Service/internal/config"
	"github.com/Chandra5468/cfp-Products-Service/internal/services/database/mongodb"
	"github.com/Chandra5468/cfp-Products-Service/internal/services/database/postgresql"
)

func main() {
	// Loading env configs
	err := config.MustLoad()

	if err != nil {
		log.Fatalf("Error while loading env file %v", err)
	}

	// Connect to psql database here
	db, err := postgresql.NewPostgres(os.Getenv("POSTGRESQL_STRING"))
	if err != nil {
		log.Fatalf("Unable to connect with postgresql database---- %v", err)
	}

	// Connect to Mongodb database here
	mdbClient, err := mongodb.NewMongodbClient(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatalf("unable to connect with mongodb database ----- %v", err)
	}

	err = mdbClient.Disconnect(context.TODO()) // pass it to httpapi.NewApiServer
	if err != nil {
		slog.Error("Alert", "message", slog.String("error : ", err.Error()))
	} else {
		slog.Info("message", "mongodb :", "disconnected sucessfully")
	}
	// Call GRPC Server here. No need to pass

	// Calling HTTP API Server
	server := httpapi.NewApiServer(os.Getenv("HTTP_ADDRESS"), db)

	server.RUN()
	// set APP_ENV=local
	//  go run .\cmd\main.go
}
