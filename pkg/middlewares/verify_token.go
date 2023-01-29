package middlewares

import (
	"context"
	"net/http"
	"strings"

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

		rawScope, exists := verifiedToken.Get("scope")
		if !exists {
			oplog.Warn().Msg("\"scope\" does not exists in the token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		strScope, ok := rawScope.(string)
		if !ok {
			oplog.Warn().Msg("\"scope\" isn't string")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		scopes := strings.Split(strScope, " ")

		c := r.Context()
		for _, scope := range scopes {
			if scope == "abura" || scope == "minmin" || scope == "kuma" || scope == "niinii" || scope == "tsukutsuku" {
				c = context.WithValue(c, scope, true)
			}
		}

		new_r := r.WithContext(c)

		next.ServeHTTP(w, new_r)
	})
}
