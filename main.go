package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

// 守るリソース ... セミに関する情報
// semis エンドポイントにアクセスすると、トークンによって違う情報がもらえる
// アブラゼミ, ミンミンゼミ、クマゼミ、ニイニイゼミ、ツクツクボウシ

const bearer_length = len("Bearer ")

func getToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header["Authorization"][0]
		token := auth[bearer_length:]

		log.Println("Incoming token:", token)

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
