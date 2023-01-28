package handlers

import (
	"encoding/json"
	"net/http"
)

type semiResource struct {
	// セミの名前
	name string
	// セミの体長
	length string
}

var semiResources = map[string]semiResource{
	"abura":      {"アブラゼミ", "5cm"},
	"minmin":     {"ミンミンゼミ", "3.5cm"},
	"kuma":       {"クマゼミ", "7cm"},
	"niinii":     {"ニイニイゼミ", "2cm"},
	"tsukutsuku": {"ツクツクボウシ", "4.5cm"},
}

func ResourcesHandler(w http.ResponseWriter, r *http.Request) {
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
}
