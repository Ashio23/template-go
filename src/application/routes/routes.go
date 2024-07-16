package routes

import (
	"net/http"
	app "template-go/src/application"
	"template-go/src/domains/example-service"
	"template-go/src/domains/health"
)

func AddRoutes(appInstance *app.App) {
	appInstance.AddRoute("GET", "/template-go/health", http.HandlerFunc(health.HealthHandler))
	appInstance.AddRoute("POST", "/template-go/api/v1/example", http.HandlerFunc(example.ExampleHandler))
}
