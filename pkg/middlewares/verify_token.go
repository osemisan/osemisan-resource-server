package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type SemiScopes struct {
	Abura      bool
	Minmin     bool
	Kuma       bool
	Niinii     bool
	Tsukutsuku bool
}

func VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		verifiedToken, err := jwt.ParseRequest(r)
		if err != nil {
			oplog.Err(err).Msg("failed to verify token from HTTP request.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		scopes, exists := verifiedToken.Get("scopes")
		if !exists {
			oplog.Warn().Msg("\"scopes\" not found from the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s, ok := scopes.(SemiScopes)
		if !ok {
			oplog.Warn().Msg("\"scopes\" isn't valid semi permissions form")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := context.WithValue(r.Context(), "abura", s.Abura)
		c = context.WithValue(c, "minmin", s.Minmin)
		c = context.WithValue(c, "kuma", s.Kuma)
		c = context.WithValue(c, "niinii", s.Niinii)
		c = context.WithValue(c, "tsukutsuku", s.Tsukutsuku)

		new_r := r.WithContext(c)

		next.ServeHTTP(w, new_r)
	})
}
