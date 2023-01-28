package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const bearer_length = len("Bearer ")

type semiResource struct {
	// セミの名前
	name string
	// セミの体長
	length string
}

var semiResources = []semiResource{
	{"アブラゼミ", "5cm"},
	{"ミンミンゼミ", "3.5cm"},
	{"クマゼミ", "7cm"},
	{"ニイニイゼミ", "2cm"},
	{"ツクツクボウシ", "4.5cm"},
}

type semiScopes struct {
	abura      bool
	minmin     bool
	kuma       bool
	niinii     bool
	tsukutsuku bool
}

func getToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())

		verifiedToken, err := jwt.ParseRequest(r)
		if err != nil {
			oplog.Err(err).Msgf("failed to verify token from HTTP request.", err)
			return
		}

		scopes, exists := verifiedToken.Get("scopes")
		if !exists {
			oplog.Warn().Msg("\"scopes\" not found from the token")
			return
		}

		s, ok := scopes.(semiScopes)
		if !ok {
			oplog.Warn().Msg("\"scopes\" isn't valid semi permissions form")
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

func main() {
	l := httplog.NewLogger("osemisan-resource-server", httplog.Options{
		JSON: true,
	})
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(l))
	r.Use(getToken)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("key:", r.Context().Value("key"))
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":9002", r)
}
