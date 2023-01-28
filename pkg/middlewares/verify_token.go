package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type semiScopes struct {
	abura      bool
	minmin     bool
	kuma       bool
	niinii     bool
	tsukutsuku bool
}

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		verifiedToken, err := jwt.ParseRequest(r)
		if err != nil {
			oplog.Err(err).Msgf("failed to verify token from HTTP request.", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		scopes, exists := verifiedToken.Get("scopes")
		if !exists {
			oplog.Warn().Msg("\"scopes\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s, ok := scopes.(semiScopes)
		if !ok {
			oplog.Warn().Msg("\"scopes\" isn't valid semi permissions form")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := context.WithValue(r.Context(), "abura", s.abura)
		c = context.WithValue(c, "minmin", s.minmin)
		c = context.WithValue(c, "kuma", s.kuma)
		c = context.WithValue(c, "niinii", s.niinii)
		c = context.WithValue(c, "tsukutsuku", s.tsukutsuku)

		new_r := r.WithContext(c)

		next.ServeHTTP(w, new_r)
	})
}
