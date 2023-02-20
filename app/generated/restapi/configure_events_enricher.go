// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swgservice "github.com/Alveona/go-base-service"

	"github.com/Alveona/go-events-enricher/app"
	"github.com/Alveona/go-events-enricher/app/generated/restapi/operations"
	"github.com/Alveona/go-events-enricher/version"
)

var (
	sc   swgservice.SwaggerConfigurator
	srv  swgservice.ServiceImplementation
	bSrv swgservice.BaseService
)

func configureFlags(api *operations.EventsEnricherAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.EventsEnricherAPI) http.Handler {
	baseSrv, err := swgservice.NewBaseService("events_enricher", version.VERSION)
	if err != nil {
		panic(fmt.Sprintf("failed to create base service: %+v", err))
	}

	srv = app.New(baseSrv)
	sc = baseSrv
	bSrv = baseSrv

	api.ServeError = sc.ServeError
	api.Logger = sc.LogCallback
	api.ServerShutdown = func() {
		sc.ShutdownCallback(srv.OnShutdown)
	}

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	if err := srv.ConfigureService(); err != nil {
		panic(fmt.Sprintf("failed to configure service: %+v", err))
	}

	srv.SetupSwaggerHandlers(api)

	customRoutes := setupCustomRoutes(api.Serve(func(handler http.Handler) http.Handler {
		return handler
	}))
	return setupMiddlewares(api.Context(), customRoutes)
}

// setupCustomRoutes creates http handler to serve custom routes
func setupCustomRoutes(next http.Handler) http.Handler {
	host, _ := os.Hostname()

	hCheck := &swgservice.Healthcheck{
		AppName:    "go-events-enricher",
		Version:    version.VERSION,
		ServerName: host,
	}
	for _, checkerFunc := range srv.HealthCheckers() {
		hCheck.AddChecker(checkerFunc)
	}

	mux := http.NewServeMux()
	mux.Handle("/", next)
	mux.Handle("/health/check", swgservice.Handler(hCheck))
	mux.Handle("/metrics", promhttp.HandlerFor(bSrv.Metrics(), promhttp.HandlerOpts{}))

	return mux
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	sc.ConfigureTLS(tlsConfig)
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
	sc.ConfigureServer(scheme, addr)
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(ctx *middleware.Context, handler http.Handler) http.Handler {
	return sc.SetupMiddlewares(ctx, handler)
}
