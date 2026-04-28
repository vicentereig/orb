package output

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestSuccessEnvelope(t *testing.T) {
	got := Success(map[string]string{"version": "dev"}, map[string]string{"operation": "version"})

	assertJSONEq(t, `{
		"success": true,
		"data": {"version": "dev"},
		"meta": {"operation": "version"},
		"error": null
	}`, got)
}

func TestErrorEnvelopeRedactsSecrets(t *testing.T) {
	got := Error(errors.New("bad Authorization bearer ORB_API_KEY"), "api_error", nil)

	assertJSONEq(t, `{
		"success": false,
		"data": null,
		"meta": null,
		"error": {"message": "bad [REDACTED] [REDACTED][REDACTED]", "type": "api_error"}
	}`, got)
}

func TestPrettyFormatsJSON(t *testing.T) {
	got, err := Pretty(`{"success":true,"data":{"ok":true},"meta":null,"error":null}`)
	if err != nil {
		t.Fatalf("Pretty returned error: %v", err)
	}
	if !json.Valid([]byte(got)) {
		t.Fatalf("Pretty returned invalid JSON: %s", got)
	}
	if got == `{"success":true,"data":{"ok":true},"meta":null,"error":null}` {
		t.Fatalf("Pretty did not format JSON")
	}
}

func assertJSONEq(t *testing.T, want, got string) {
	t.Helper()

	var wantValue interface{}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("want is invalid JSON: %v", err)
	}
	var gotValue interface{}
	if err := json.Unmarshal([]byte(got), &gotValue); err != nil {
		t.Fatalf("got is invalid JSON: %v\n%s", err, got)
	}

	wantBytes, _ := json.Marshal(wantValue)
	gotBytes, _ := json.Marshal(gotValue)
	if string(wantBytes) != string(gotBytes) {
		t.Fatalf("JSON mismatch\nwant: %s\n got: %s", wantBytes, gotBytes)
	}
}
