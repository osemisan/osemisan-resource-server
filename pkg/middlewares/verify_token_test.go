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

type scopes struct {
	abura      bool
	minmin     bool
	kuma       bool
	niinii     bool
	tsukutsuku bool
}

func GetTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}
	return fn
}

func BuildSimpleJwt(t *testing.T) string {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Build()
	if err != nil {
		t.Error("Failed to build JWT", err)
	}
	return fmt.Sprintf("Bearer %s", tok)
}

func BuildScopedJwt(t *testing.T, s struct {
	abura      bool
	minmin     bool
	kuma       bool
	niinii     bool
	tsukutsuku bool
}) string {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Claim("scopeAbura", s.abura).
		Claim("scopeMinmin", s.minmin).
		Claim("scopeKuma", s.kuma).
		Claim("scopeNiinii", s.niinii).
		Claim("scopeTsukutsuku", s.tsukutsuku).
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
		wantScopes     scopes
	}{
		{
			"リクエストヘッダからJWTが読み出せなかった401",
			"invalid-token",
			http.StatusUnauthorized,
			scopes{},
		},
		{
			"JWTを読み解くことはできたが、scopesが含まれていないとき401",
			BuildSimpleJwt(t),
			http.StatusUnauthorized,
			scopes{},
		},
		{
			"アブラゼミだけ閲覧可能なスコープが付与されている",
			BuildScopedJwt(t, scopes{true, false, false, false, false}),
			http.StatusAccepted,
			scopes{true, false, false, false, false},
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
				t.Errorf("Unexpected status code, expected: %d, actural: %d", tt.wantStatusCode, res.StatusCode)
			}

			if v, ok := req.Context().Value("abura").(bool); ok {
				if v != tt.wantScopes.abura {
					t.Error("Unexpected scopes.Abura")
				}
			}
			if v, ok := req.Context().Value("minmin").(bool); ok {
				if v != tt.wantScopes.minmin {
					t.Error("Unexpected scopes.Minmin")
				}
			}
			if v, ok := req.Context().Value("kuma").(bool); ok {
				if v != tt.wantScopes.kuma {
					t.Error("Unexpected scopes.Kuma")
				}
			}
			if v, ok := req.Context().Value("Niinii").(bool); ok {
				if v != tt.wantScopes.niinii {
					t.Error("Unexpected scopes.Niinii")
				}
			}
			if v, ok := req.Context().Value("tsukutsuku").(bool); ok {
				if v != tt.wantScopes.tsukutsuku {
					t.Error("Unexpected scopes.Tsukutsuku")
				}
			}
		})
	}
}
