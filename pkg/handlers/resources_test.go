package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osemisan/osemisan-resource-server/pkg/handlers"
	"github.com/osemisan/osemisan-resource-server/pkg/middlewares"
	"github.com/osemisan/osemisan-resource-server/testutil"
)

func GetTestHandler() http.HandlerFunc {
	fn := handlers.ResourcesHandler
	return fn
}

func TestResourcesHandler(t *testing.T) {
	s := httptest.NewServer(middlewares.VerifyToken(GetTestHandler()))
	defer s.Close()

	c := new(http.Client)

	tests := []struct {
		name           string
		token          string
		wantStatusCode int
		wantJson       string
	}{
		{
			"アブラゼミのみが閲覧できるトークに対してアブラゼミだけのレスポンス",
			testutil.BuildScopedJwt(t, testutil.Scopes{Abura: true, Minmin: false, Kuma: false, Niinii: false, Tsukutsuku: false}),
			http.StatusOK,
			`[{"name":"アブラゼミ","length":"5cm"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, s.URL, nil)
			if err != nil {
				t.Error("Failed to create reeust", err)
				return
			}

			req.Header.Set("Authorization", tt.token)

			res, err := c.Do(req)
			if err != nil {
				t.Error("Failed to request", err)
				return
			}
			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("Unexpected status code, expected: %d, actual: %d", tt.wantStatusCode, res.StatusCode)
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				t.Error("Failed to read response body")
				return
			}

			if string(b) != tt.wantJson {
				t.Errorf("Unexpected response, expected: %s, actual: %s", tt.wantJson, string(b))
			}
		})
	}
}
