package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

// 守るリソース ... セミに関する情報
// semis エンドポイントにアクセスすると、トークンによって違う情報がもらえる
// アブラゼミ, ミンミンゼミ、クマゼミ、ニイニイゼミ、ツクツクボウシ

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

func getToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		auth := r.Header["Authorization"][0]
		token := auth[bearer_length:]

		oplog.Info().Msgf("Incoming token: %s", token)

		new_r := r.WithContext(context.WithValue(r.Context(), "key", "value"))
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
	http.ListenAndServe(":3000", r)
}
