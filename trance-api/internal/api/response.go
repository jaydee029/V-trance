package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type errresponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, res string) {
	if code > 499 {
		log.Printf("Responding with %d error: %s", code, res)
	}
	type errresponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errresponse{
		Error: res,
	})
}

func respondWithJson(w http.ResponseWriter, code int, res interface{}) {
	w.Header().Set("content-type", "application/json")
	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
