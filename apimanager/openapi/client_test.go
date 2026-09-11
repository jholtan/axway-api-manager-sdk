package openapi

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthClientAppliesCredentialsToGeneratedOperation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/alerts" {
			t.Errorf("path = %q, want /alerts", r.URL.Path)
		}
		username, password, ok := r.BasicAuth()
		if !ok || username != "user" || password != "password" {
			t.Errorf("basic auth = (%q, %q, %t), want (user, password, true)", username, password, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewBasicAuthClient(server.URL, "user", "password", server.Client())
	if err != nil {
		t.Errorf("NewBasicAuthClient() error = %v", err)
	}
	resp, err := client.GetAlerts(context.Background())
	if err != nil {
		t.Errorf("GetAlerts() error = %v", err)
	}
	resp.Body.Close()
}

func TestGeneratedClientSupportsRawRequestBodies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/apirepo/import" {
			t.Errorf("path = %q, want /apirepo/import", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ReadAll() error = %v", err)
		}
		if !bytes.Equal(body, []byte("definition")) {
			t.Errorf("body = %q, want definition", body)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client, err := NewBasicAuthClient(server.URL, "user", "password", server.Client())
	if err != nil {
		t.Errorf("NewBasicAuthClient() error = %v", err)
	}
	resp, err := client.ImportApisFromFile2WithBody(context.Background(), "application/octet-stream", bytes.NewReader([]byte("definition")))
	if err != nil {
		t.Errorf("ImportApisFromFile2WithBody() error = %v", err)
	}
	resp.Body.Close()
}

func TestDoReturnsRawSuccessfulResponsesAndRejectsErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/failure" {
			http.Error(w, "invalid", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("image"))
	}))
	defer server.Close()

	client, err := NewBasicAuthClient(server.URL, "user", "password", server.Client())
	if err != nil {
		t.Errorf("NewBasicAuthClient() error = %v", err)
	}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/success", nil)
	if err != nil {
		t.Errorf("NewRequestWithContext() error = %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Do() error = %v", err)
	}
	body, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil || string(body) != "image" {
		t.Errorf("response body = %q, error = %v", body, err)
	}

	request, err = http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL+"/failure", nil)
	if err != nil {
		t.Errorf("NewRequestWithContext() error = %v", err)
	}
	if _, err := client.Do(request); err == nil {
		t.Error("Do() error = nil, want non-success error")
	}
}
