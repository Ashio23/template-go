package backend

type ExampleResponseGenericError struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"error"`
}

func (e *ExampleResponseGenericError) Error() string {
	return e.ErrorMessage
}

type ExampleResponse404Error struct {
	Timestamp    string `json:"timestamp"`
	Status       int    `json:"status"`
	ErrorMessage string `json:"error"`
	Path         string `json:"path"`
}

func (e *ExampleResponse404Error) Error() string {
	return e.ErrorMessage
}
