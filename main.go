package main

import (
	"context"
	"encoding/json"
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

var semiResources = map[string]semiResource{
	"abura": {"アブラゼミ", "5cm"},
	"minmin": {"ミンミンゼミ", "3.5cm"},
	"kuma": {"クマゼミ", "7cm"},
	"niinii": {"ニイニイゼミ", "2cm"},
	"tsukutsuku": {"ツクツクボウシ", "4.5cm"},
}

type semiScopes struct {
	abura      bool
	minmin     bool
	kuma       bool
	niinii     bool
	tsukutsuku bool
}

func verifyToken(next http.Handler) http.Handler {
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

func main() {
	l := httplog.NewLogger("osemisan-resource-server", httplog.Options{
		JSON: true,
	})
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(l))
	r.Use(verifyToken)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/resources", func(w http.ResponseWriter, r *http.Request) {
		resources := make([]semiResource, 0, 5)

		if r.Context().Value("abura").(bool) {
			resources = append(resources, semiResources["abura"])
		}
		if r.Context().Value("minmin").(bool) {
			resources = append(resources, semiResources["minmin"])
		}
		if r.Context().Value("kuma").(bool) {
			resources = append(resources, semiResources["kuma"])
		}
		if r.Context().Value("niinii").(bool) {
			resources = append(resources, semiResources["niinii"])
		}
		if r.Context().Value("tsukutsuku").(bool) {
			resources = append(resources, semiResources["tsukutsuku"])
		}

		jsonRes, err := json.Marshal(resources)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRes)
	})

	http.ListenAndServe(":9002", r)
}
