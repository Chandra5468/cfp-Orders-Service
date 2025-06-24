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

	v1 "github.com/Chandra5468/cfp-Products-Service/internal/handlers/http/v1"
	"github.com/Chandra5468/cfp-Products-Service/internal/middleware"
	"github.com/Chandra5468/cfp-Products-Service/internal/services/database/postgresql/orders"
	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	addr string
	db   *sql.DB
	// mdb In future if we want to migrate for unstructred database use mongo conn here
}

func NewApiServer(httpAddress string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: httpAddress,
		db:   db,
	}
}

func (a *APIServer) RUN() {
	// router := http.NewServeMux() // Default inbuilt http router

	router := chi.NewRouter()

	// Currently directing all calls for v1 version.

	// router.Use(middleware.Logger
	newHandler := middleware.CorsHandler(router)
	ordersStore := orders.NewStore(a.db)        // call services first (from postgresql database)
	ordersHandler := v1.NewHandler(ordersStore) // assign those services to handler. Internally implements interfaces. So, if services come from mongo in future they need to implement those services
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
	err2 := server.Shutdown(ctx)

	if err2 != nil {
		slog.Error("failed to shutdown server", slog.String("error", err2.Error()))
	} else {
		slog.Info("Server shutdown successful")
	}

}
