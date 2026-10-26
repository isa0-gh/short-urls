package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterRoutesServesPagesAndAssets(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	tests := []struct {
		path        string
		wantStatus  int
		wantContent string
	}{
		{path: "/login", wantStatus: http.StatusOK, wantContent: `data-page="login"`},
		{path: "/register", wantStatus: http.StatusOK, wantContent: `data-page="register"`},
		{path: "/dashboard", wantStatus: http.StatusOK, wantContent: `data-page="dashboard"`},
		{path: "/assets/app.css", wantStatus: http.StatusOK, wantContent: ".navbar"},
		{path: "/assets/app.js", wantStatus: http.StatusOK, wantContent: "setupDashboard"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, test.path, nil))

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if !strings.Contains(response.Body.String(), test.wantContent) {
				t.Fatalf("response body does not contain %q", test.wantContent)
			}
		})
	}
}
