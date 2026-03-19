package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Request struct {
	Method      string
	URL         string
	Body        []byte
	Headers     map[string]string
	QueryParams map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

func Do(req *Request) (*Response, error) {
	// Build URL with query params
	url := req.URL
	if len(req.QueryParams) > 0 {
		separator := "?"
		if strings.Contains(url, "?") {
			separator = "&"
		}
		var params []string
		for k, v := range req.QueryParams {
			params = append(params, k+"="+v)
		}
		url += separator + strings.Join(params, "&")
	}

	var bodyReader io.Reader
	if len(req.Body) > 0 {
		bodyReader = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(strings.ToUpper(req.Method), url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if len(req.Body) > 0 && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// PrintResponse prints the response to stdout with pretty JSON formatting.
func PrintResponse(resp *Response) {
	// Try to pretty-print JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, resp.Body, "", "  "); err == nil {
		fmt.Fprintln(os.Stdout, prettyJSON.String())
	} else {
		// Not JSON, print raw
		fmt.Fprintln(os.Stdout, string(resp.Body))
	}
}

// PrintResponseWithStatus prints status code to stderr and body to stdout.
func PrintResponseWithStatus(resp *Response) {
	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP %d\n", resp.StatusCode)
	}
	PrintResponse(resp)
}
