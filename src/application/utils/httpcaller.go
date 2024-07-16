package utils

import (
	"bytes"
	"fmt"
	"net/http"
)

type HttpCaller struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}

func (h *HttpCaller) Call() (*http.Response, error) {
	if h.Body == nil {
		h.Body = []byte{}
	}
	client := &http.Client{}
	req, err := http.NewRequest(h.Method, h.Url, bytes.NewBuffer(h.Body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if h.Headers != nil {
		for k, v := range h.Headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling %s: %v", h.Url, err)
	}

	return resp, nil
}
