package httpapi

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcclient "github.com/Chandra5468/cfp-Products-Service/cmd/grpcClient"
	v1 "github.com/Chandra5468/cfp-Products-Service/internal/handlers/http/v1"
	"github.com/Chandra5468/cfp-Products-Service/internal/middleware"
	mdbDatastore "github.com/Chandra5468/cfp-Products-Service/internal/services/database/mongodb/orders"
	psqlDatastore "github.com/Chandra5468/cfp-Products-Service/internal/services/database/postgresql/orders"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type APIServer struct {
	addr string
	db   *sql.DB
	mdb  *mongo.Client
}

func NewApiServer(httpAddress string, db *sql.DB, mdbClient *mongo.Client) *APIServer {
	return &APIServer{
		addr: httpAddress,
		db:   db,
		mdb:  mdbClient,
	}
}

func (a *APIServer) RUN() {
	// router := http.NewServeMux() // Default inbuilt http router

	router := chi.NewRouter()

	// Currently directing all calls for v1 version.

	// router.Use(middleware.Logger
	newHandler := middleware.CorsHandler(router)

	ordersStore := psqlDatastore.NewStore(a.db) // call services first (from postgresql database)
	complaintsStore := mdbDatastore.NewStore(a.mdb)
	// Grpc Conn
	conn := grpcclient.NewGrpcClient("localhost:9002")
	ordersHandler := v1.NewHandler(ordersStore, complaintsStore, conn) // assign those services to handler. Internally implements interfaces. So, if services come from mongo in future they need to implement those services
	ordersHandler.RegisterRoutes(router)

	server := &http.Server{
		Addr:    a.addr,
		Handler: newHandler,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		slog.Info("message", "Listening on address", a.addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", slog.String("error", err.Error()))
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	slog.Error("Alert ", "message", "Shutting down http server ")

	slog.Info("message", "closing", "postgresql")
	err := a.db.Close()
	if err != nil {
		slog.Error("Error closing PostgreSQL", slog.String("error", err.Error()))
	} else {
		slog.Info("PostgreSQL database closed successfully")
	}
	slog.Info("message", "closing", "mongodatabase")
	err = a.mdb.Disconnect(context.TODO())
	if err != nil {
		slog.Error("Error closing MongoDB", slog.String("error", err.Error()))
	} else {
		slog.Info("Mongo database closed successfully")
	}
	err2 := server.Shutdown(ctx)

	if err2 != nil {
		slog.Error("failed to shutdown server", slog.String("error", err2.Error()))
	} else {
		slog.Info("Server shutdown successful")
	}

}
