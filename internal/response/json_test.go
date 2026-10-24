package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONNilWritesEmptyBody(t *testing.T) {
	recorder := httptest.NewRecorder()

	if err := JSON(recorder, http.StatusOK, nil); err != nil {
		t.Fatalf("JSON returned an error: %v", err)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", recorder.Body.String())
	}
}
