package api

import (
	"encoding/json"
	"net/http"

	"go.mattglei.ch/go.mattglei.ch/internal/util"
	"go.mattglei.ch/timber"
)

var response []byte

func init() {
	data, err := json.Marshal(struct {
		Ok bool `json:"ok"`
	}{Ok: true})
	if err != nil {
		timber.Fatal(err, "failed to set health check response")
		return
	}
	response = data
}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(response)
	if err != nil {
		util.InternalServerError(w, err)
		return
	}
}
