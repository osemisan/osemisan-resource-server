package testutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const symKey = "hoge"

func BuildSimpleJwt(t *testing.T) string {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Build()
	if err != nil {
		t.Error("Failed to build JWT", err)
	}
	key, err := jwk.FromRaw([]byte(symKey))
	if err != nil {
		t.Error("Failed to create key from raw", err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		t.Error("Failed to sign token", err)
	}
	return fmt.Sprintf("Bearer %s", signed)
}

type Scopes struct {
	Abura      bool
	Minmin     bool
	Kuma       bool
	Niinii     bool
	Tsukutsuku bool
}

func BuildScopedJwt(t *testing.T, s Scopes) string {
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Claim("scopeAbura", s.Abura).
		Claim("scopeMinmin", s.Minmin).
		Claim("scopeKuma", s.Kuma).
		Claim("scopeNiinii", s.Niinii).
		Claim("scopeTsukutsuku", s.Tsukutsuku).
		Build()
	if err != nil {
		t.Error("Failed to build JWT", err)
	}
	key, err := jwk.FromRaw([]byte("hoge"))
	if err != nil {
		t.Error("Failed to create key from raw", err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		t.Error("Failed to sign token", err)
	}
	return fmt.Sprintf("Bearer %s", signed)
}
