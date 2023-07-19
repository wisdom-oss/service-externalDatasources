package main

import (
	managementRoutes "external-api-service/routes/management"
	"fmt"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v2"
	"net/http"
	"os"
	"os/signal"
	"time"

	"external-api-service/globals"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.With().Logger()
	l.Info().Msgf("starting %s service", globals.ServiceName)

	// create a new router
	router := chi.NewRouter()
	// add some middlewares to the router to allow identifying requests
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(httplog.Handler(l))
	// now add the authorization middleware to the router
	router.Use(wisdomMiddleware.Authorization(globals.AuthorizationConfiguration, globals.ServiceName))
	router.Use(wisdomMiddleware.NativeErrorHandler(globals.ServiceName))
	router.Use(wisdomMiddleware.WISdoMErrorHandler(globals.Errors))
	// now mount the admin router
	router.Mount("/management", managementRouter())

	// now boot up the service
	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", globals.Environment["LISTEN_PORT"]),
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      router,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); err != nil {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// Block further code execution until the shutdown signal was received
	<-cancelSignal

}

func managementRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", managementRoutes.AllExternalDataSources)
	r.Post("/", managementRoutes.NewDataSource)
	r.Get("/{dataSourceUUID}", managementRoutes.SingleDataSource)
	r.Put("/{dataSourceUUID}", managementRoutes.ReplaceDataSourceRepresentation)
	r.Patch("/{dataSourceUUID}", managementRoutes.UpdateDataSourceRepresentation)
	r.Delete("/{dataSourceUUID}", managementRoutes.DeleteExternalDataSource)
	return r
}
