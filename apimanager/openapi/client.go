package openapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func NewBasicAuthClient(server, username, password string, httpClient *http.Client) (*Client, error) {
	options := []ClientOption{
		WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.SetBasicAuth(username, password)
			return nil
		}),
	}
	if httpClient != nil {
		options = append(options, WithHTTPClient(httpClient))
	}
	return NewClient(server, options...)
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	for _, editor := range c.RequestEditors {
		if err := editor(req.Context(), req); err != nil {
			return nil, err
		}
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return resp, nil
	}

	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
}
