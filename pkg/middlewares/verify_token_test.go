package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/osemisan/osemisan-resource-server/pkg/middlewares"
)

func GetTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	return fn
}

func BuildJwt(t *testing.T) string {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Build()
	if err != nil {
		t.Error("Failed to build JWT", err)
	}
	return fmt.Sprintf("Bearer %s", tok)
}

func TestVerifyToken(t *testing.T) {
	s := httptest.NewServer(middlewares.VerifyToken(GetTestHandler()))
	defer s.Close()

	c := new(http.Client)

	tests := []struct {
		name           string
		token          string
		wantStatusCode int
		wantScopes middlewares.SemiScopes
	}{
		{
			"リクエストヘッダからJWTが読み出せなかった401",
			"invalid-token",
			http.StatusUnauthorized,
			middlewares.SemiScopes{},
		},
		{
			"JWTを読み解くことはできたが、scopesが含まれていないとき401",
			BuildJwt(t),
			http.StatusUnauthorized,
			middlewares.SemiScopes{},
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
				t.Error("Unexpected status code")
			}

			if v, ok := req.Context().Value("abura").(bool); ok {
				if v != tt.wantScopes.Abura {
					t.Error("Unexpected scopes.Abura")
				}
			}
			if v, ok := req.Context().Value("minmin").(bool); ok {
				if v != tt.wantScopes.Minmin {
					t.Error("Unexpected scopes.Minmin")
				}
			}
			if v, ok := req.Context().Value("kuma").(bool); ok {
				if v != tt.wantScopes.Kuma {
					t.Error("Unexpected scopes.Kuma")
				}
			}
			if v, ok := req.Context().Value("Niinii").(bool); ok {
				if v != tt.wantScopes.Niinii {
					t.Error("Unexpected scopes.Niinii")
				}
			}
			if v, ok := req.Context().Value("tsukutsuku").(bool); ok {
				if v != tt.wantScopes.Tsukutsuku {
					t.Error("Unexpected scopes.Tsukutsuku")
				}
			}
		})
	}
}
