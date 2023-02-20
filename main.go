package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/loads"
	"github.com/vrischmann/envconfig"
	_ "go.uber.org/automaxprocs"

	"github.com/Alveona/go-events-enricher/app/generated/restapi"
	"github.com/Alveona/go-events-enricher/app/generated/restapi/operations"
	"github.com/Alveona/go-events-enricher/version"
)

func main() {

	fmt.Printf("BuildTime - %s, Version - %s\n", version.BUILD_TIME, version.VERSION)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewEventsEnricherAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer func() {
		if err := server.Shutdown(); err != nil {
			//nolint:gocritic
			log.Fatalln(err) // disable gocritic due to os.Exit(1) in Fatatlln and exitAfterDefer
		}
	}()

	server.ConfigureFlags()
	server.ConfigureAPI()

	var conf struct {
		HTTPBindPort    int
		GracefulTimeout time.Duration `envconfig:"default=1m"`
	}
	if err := envconfig.InitWithPrefix(&conf, "events_enricher"); err != nil {
		//nolint:gocritic
		log.Fatalln(err) // disable gocritic due to os.Exit(1) in Fatatlln and exitAfterDefer
	}
	server.Port = conf.HTTPBindPort
	server.GracefulTimeout = conf.GracefulTimeout
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
