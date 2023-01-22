package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// 守るリソース ... セミに関する情報
// semis エンドポイントにアクセスすると、トークンによって違う情報がもらえる
// アブラゼミ, ミンミンゼミ、クマゼミ、ニイニイゼミ、ツクツクボウシ

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		new_r := r.WithContext(context.WithValue(r.Context(), "key", "value"))
		next.ServeHTTP(w, new_r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger, myMiddleware)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("key:", r.Context().Value("key"))
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}
