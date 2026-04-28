package output

import (
	"encoding/json"
	"fmt"
)

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Result struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data"`
	Meta    interface{}  `json:"meta"`
	Error   *ErrorDetail `json:"error"`
}

func Success(data interface{}, meta interface{}) string {
	return marshal(Result{
		Success: true,
		Data:    data,
		Meta:    meta,
		Error:   nil,
	})
}

func Error(err error, typ string, meta interface{}) string {
	if typ == "" {
		typ = "error"
	}
	return marshal(Result{
		Success: false,
		Data:    nil,
		Meta:    meta,
		Error: &ErrorDetail{
			Message: redact(err.Error()),
			Type:    typ,
		},
	})
}

func marshal(result Result) string {
	b, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf(`{"success":false,"data":null,"meta":null,"error":{"message":"%s","type":"marshal_error"}}`, redact(err.Error()))
	}
	return string(b)
}

func Pretty(jsonLine string) (string, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(jsonLine), &v); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func redact(s string) string {
	secrets := []string{"Authorization", "ORB_API_KEY", "api_key", "bearer "}
	out := s
	for _, secret := range secrets {
		if len(out) == 0 {
			return out
		}
		out = replaceCaseInsensitive(out, secret, "[REDACTED]")
	}
	return out
}

func replaceCaseInsensitive(s, old, replacement string) string {
	if old == "" {
		return s
	}
	b := []byte(s)
	lower := []byte(s)
	target := []byte(old)
	for i := range lower {
		if lower[i] >= 'A' && lower[i] <= 'Z' {
			lower[i] += 'a' - 'A'
		}
	}
	for i := range target {
		if target[i] >= 'A' && target[i] <= 'Z' {
			target[i] += 'a' - 'A'
		}
	}
	for i := 0; i+len(target) <= len(lower); i++ {
		match := true
		for j := range target {
			if lower[i+j] != target[j] {
				match = false
				break
			}
		}
		if match {
			return string(b[:i]) + replacement + string(b[i+len(target):])
		}
	}
	return s
}
