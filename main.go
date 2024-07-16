package main

import (
	"net/http"
	"os"
	app "template-go/src/application"
	"template-go/src/application/middlewares"
	"template-go/src/application/routes"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	log.TextFormatter
}

func init() {
	if os.Getenv("ENV") != "" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&CustomFormatter{
			TextFormatter: log.TextFormatter{
				DisableTimestamp: true,
				ForceColors:      true,
			},
		})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appInstance := &app.App{Mux: http.NewServeMux()}
	routes.AddRoutes(appInstance)

	for _, route := range appInstance.Routes {
		log.Infof("Mapped route: %s %s", route.Method, route.Path)
	}

	loggedRouter := middlewares.LoggingMiddleware()(appInstance.Mux)

	env := os.Getenv("ENV")
	if env != "local" {
		lambda.Start(httpadapter.New(loggedRouter).ProxyWithContext)
		return
	}

	addr := os.Getenv("HOST") + ":" + os.Getenv("PORT")
	log.Infof("Server is running on %s", addr)
	if err := http.ListenAndServe(addr, loggedRouter); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
