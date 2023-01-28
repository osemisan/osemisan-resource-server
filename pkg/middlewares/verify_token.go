package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const symKey = "hoge"

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		key, err := jwk.FromRaw([]byte(symKey))
		if err != nil {
			oplog.Err(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		verifiedToken, err := jwt.ParseRequest(r, jwt.WithKey(jwa.HS256, key))
		if err != nil {
			oplog.Err(err).Msg("failed to verify token from HTTP request.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sAbura, exists := verifiedToken.Get("scopeAbura")
		if !exists {
			oplog.Warn().Msg("\"scopeAbura\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sMinmin, exists := verifiedToken.Get("scopeMinmin")
		if !exists {
			oplog.Warn().Msg("\"scopeMinmin\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sKuma, exists := verifiedToken.Get("scopeKuma")
		if !exists {
			oplog.Warn().Msg("\"scopeKuma\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sNiinii, exists := verifiedToken.Get("scopeNiinii")
		if !exists {
			oplog.Warn().Msg("\"scopeNiinii\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sTsukutsuku, exists := verifiedToken.Get("scopeTsukutsuku")
		if !exists {
			oplog.Warn().Msg("\"scopeTsukutsuku\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := context.WithValue(r.Context(), "abura", sAbura)
		c = context.WithValue(c, "minmin", sMinmin)
		c = context.WithValue(c, "kuma", sKuma)
		c = context.WithValue(c, "niinii", sNiinii)
		c = context.WithValue(c, "tsukutsuku", sTsukutsuku)

		new_r := r.WithContext(c)

		next.ServeHTTP(w, new_r)
	})
}
