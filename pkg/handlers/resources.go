package handlers

import (
	"encoding/json"
	"net/http"
)

type semiResource struct {
	// セミの名前
	Name string `json:"name"`
	// セミの体長
	Length string `json:"length"`
}

var semiResources = map[string]semiResource{
	"abura":      {Name: "アブラゼミ", Length: "5cm"},
	"minmin":     {Name: "ミンミンゼミ", Length: "3.5cm"},
	"kuma":       {Name: "クマゼミ", Length: "7cm"},
	"niinii":     {Name: "ニイニイゼミ", Length: "2cm"},
	"tsukutsuku": {Name: "ツクツクボウシ", Length: "4.5cm"},
}

func ResourcesHandler(w http.ResponseWriter, r *http.Request) {
	resources := make([]semiResource, 0, 5)

	if v, ok := r.Context().Value("abura").(bool); ok {
		if v {
			resources = append(resources, semiResources["abura"])
		}
	}
	if v, ok := r.Context().Value("minmin").(bool); ok {
		if v {
			resources = append(resources, semiResources["minmin"])
		}
	}
	if v, ok := r.Context().Value("kuma").(bool); ok {
		if v {
			resources = append(resources, semiResources["kuma"])
		}
	}
	if v, ok := r.Context().Value("niinii").(bool); ok {
		if v {
			resources = append(resources, semiResources["niinii"])
		}
	}
	if v, ok := r.Context().Value("tsukutsuku").(bool); ok {
		if v {
			resources = append(resources, semiResources["tsukutsuku"])
		}
	}

	jsonRes, err := json.Marshal(resources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
