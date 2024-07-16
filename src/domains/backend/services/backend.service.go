package backendservices

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"template-go/src/application/utils"
	exampleback "template-go/src/domains/backend/models"
)

func ExampleServiceBack(exampleRequest exampleback.ExampleRequest) (any, error) {
	hostname := os.Getenv("EXAMPLE_SERVICE_URL")

	path := "/controller/service"

	body, err := json.Marshal(exampleRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	h := utils.HttpCaller{
		Method: "POST",
		Url:    hostname + path,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	log.Println("Calling external service:", h.Url)
	log.Println("Request body:", string(body))

	resp, err := h.Call()
	if err != nil {
		fmt.Println("Error calling external service:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil, nil
}
