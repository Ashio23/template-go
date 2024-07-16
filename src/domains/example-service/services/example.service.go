package exampleservices

import (
	example "template-go/src/domains/example-service/models"
)

func ExampleService(exampleRequest example.ExampleRequest) (*example.ExampleResponse, error) {

	data := exampleRequest
	response := &example.ExampleResponse{
		Id: data.Id,
	}

	return response, nil
}
