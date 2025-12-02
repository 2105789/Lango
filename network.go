package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HttpResponse represents the response from an HTTP request
type HttpResponse struct {
	Status     int
	StatusText string
	Headers    map[string]interface{}
	Body       string
}

// NativeHttpGet is a built-in function for HTTP GET requests
type NativeHttpGet struct{}

func (n *NativeHttpGet) Arity() int {
	return 1 // url (headers optional)
}

func (n *NativeHttpGet) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) < 1 {
		return nil, fmt.Errorf("httpGet requires at least 1 argument (url)")
	}

	url, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("httpGet: url must be a string")
	}

	var headers map[string]interface{}
	if len(arguments) > 1 {
		headers, ok = arguments[1].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("httpGet: headers must be an object")
		}
	}

	return makeHttpRequest("GET", url, "", headers)
}

func (n *NativeHttpGet) String() string {
	return "<native fn httpGet>"
}

// NativeHttpPost is a built-in function for HTTP POST requests
type NativeHttpPost struct{}

func (n *NativeHttpPost) Arity() int {
	return 2 // url, body (headers optional)
}

func (n *NativeHttpPost) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		return nil, fmt.Errorf("httpPost requires at least 2 arguments (url, body)")
	}

	url, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("httpPost: url must be a string")
	}

	body, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("httpPost: body must be a string")
	}

	var headers map[string]interface{}
	if len(arguments) > 2 {
		headers, ok = arguments[2].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("httpPost: headers must be an object")
		}
	}

	return makeHttpRequest("POST", url, body, headers)
}

func (n *NativeHttpPost) String() string {
	return "<native fn httpPost>"
}

// NativeHttpPut is a built-in function for HTTP PUT requests
type NativeHttpPut struct{}

func (n *NativeHttpPut) Arity() int {
	return 2 // url, body (headers optional)
}

func (n *NativeHttpPut) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) < 2 {
		return nil, fmt.Errorf("httpPut requires at least 2 arguments (url, body)")
	}

	url, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("httpPut: url must be a string")
	}

	body, ok := arguments[1].(string)
	if !ok {
		return nil, fmt.Errorf("httpPut: body must be a string")
	}

	var headers map[string]interface{}
	if len(arguments) > 2 {
		headers, ok = arguments[2].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("httpPut: headers must be an object")
		}
	}

	return makeHttpRequest("PUT", url, body, headers)
}

func (n *NativeHttpPut) String() string {
	return "<native fn httpPut>"
}

// NativeHttpDelete is a built-in function for HTTP DELETE requests
type NativeHttpDelete struct{}

func (n *NativeHttpDelete) Arity() int {
	return 1 // url (headers optional)
}

func (n *NativeHttpDelete) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	if len(arguments) < 1 {
		return nil, fmt.Errorf("httpDelete requires at least 1 argument (url)")
	}

	url, ok := arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("httpDelete: url must be a string")
	}

	var headers map[string]interface{}
	if len(arguments) > 1 {
		headers, ok = arguments[1].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("httpDelete: headers must be an object")
		}
	}

	return makeHttpRequest("DELETE", url, "", headers)
}

func (n *NativeHttpDelete) String() string {
	return "<native fn httpDelete>"
}

// makeHttpRequest performs the actual HTTP request
func makeHttpRequest(method, url, body string, headers map[string]interface{}) (*LangoInstance, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var reqBody io.Reader
	if body != "" {
		reqBody = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set default Content-Type for POST/PUT if not specified
	if (method == "POST" || method == "PUT") && body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	if headers != nil {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Create response object
	responseClass := &LangoClass{
		name:    "HttpResponse",
		methods: make(map[string]*LangoFunction),
	}

	instance := NewLangoInstance(responseClass)
	instance.fields["status"] = float64(resp.StatusCode)
	instance.fields["statusText"] = resp.Status
	instance.fields["body"] = string(bodyBytes)

	// Convert headers to Lango map
	headersMap := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) > 0 {
			headersMap[key] = values[0]
		}
	}
	instance.fields["headers"] = headersMap

	return instance, nil
}
