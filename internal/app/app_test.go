package app

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type fakeClient struct {
	pingResponse interface{}
	pingErr      error
}

func (f fakeClient) Ping(ctx context.Context) (interface{}, error) {
	if f.pingErr != nil {
		return nil, f.pingErr
	}
	return f.pingResponse, nil
}

func TestVersion(t *testing.T) {
	result := New(fakeClient{}, "v1.2.3").Version()

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data.Version != "v1.2.3" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestPingSuccess(t *testing.T) {
	result := New(fakeClient{pingResponse: map[string]string{"response": "pong"}}, "dev").Ping(context.Background())

	var payload struct {
		Success bool `json:"success"`
		Data    struct {
			Response string `json:"response"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !payload.Success || payload.Data.Response != "pong" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestPingError(t *testing.T) {
	result := New(fakeClient{pingErr: errors.New("unauthorized")}, "dev").Ping(context.Background())

	var payload struct {
		Success bool `json:"success"`
		Error   struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if payload.Success || payload.Error.Message != "unauthorized" || payload.Error.Type != "api_error" {
		t.Fatalf("unexpected result: %s", result)
	}
}
