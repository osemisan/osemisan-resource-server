package testutil

import (
	"fmt"
	"strings"
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
	scopes := make([]string, 5, 5)
	if s.Abura {
		scopes = append(scopes, "abura")
	}
	if s.Minmin {
		scopes = append(scopes, "minmin")
	}
	if s.Kuma {
		scopes = append(scopes, "kuma")
	}
	if s.Niinii {
		scopes = append(scopes, "niinii")
	}
	if s.Tsukutsuku {
		scopes = append(scopes, "tsukutsuku")
	}
	tok, err := jwt.NewBuilder().
		Issuer(`github.com/osemisan/osemisan-resource-server`).
		IssuedAt(time.Now()).
		Claim("scope", strings.Join(scopes, " ")).
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
