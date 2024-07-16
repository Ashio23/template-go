package example

import (
	"net/http"
	"template-go/src/application/utils"
	"template-go/src/application/validation"
	"template-go/src/domains/backend"
	example "template-go/src/domains/example-service/models"
	exampleservices "template-go/src/domains/example-service/services"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {

	exampleRequest := example.ExampleRequest{}
	if err := validation.Validate(r, &exampleRequest); err != nil {
		validation.HandleValidationError(w, err)
		return
	}
	data, err := exampleservices.ExampleService(exampleRequest)
	if err != nil {
		ErrorHandler(w, err)
		return
	}

	utils.XMLResponse(w, data, http.StatusOK)
}

func ErrorHandler(w http.ResponseWriter, err error) {
	switch err := err.(type) {
	case *backend.ExampleResponse404Error:
		utils.XMLResponse(w, err, err.Status)
	default:
		errorData := validation.ApiError{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: err.Error(),
		}
		utils.XMLResponse(w, errorData, errorData.StatusCode)
	}
}
