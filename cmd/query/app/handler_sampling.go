package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	ui "github.com/jaegertracing/jaeger/model/json"
)

func (aH *APIHandler) getAllSamplings(w http.ResponseWriter, r *http.Request) {
	retMe, err := aH.samplingWriter.GetAllSamplings()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retMe)
}

func (aH *APIHandler) getSampling(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	samp, err := aH.samplingWriter.GetSampling(vars[samplingParam])
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(samp)
}

func (aH *APIHandler) deleteSampling(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := aH.samplingWriter.DeleteSampling(vars[samplingParam])
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}
}

func (aH *APIHandler) writeSampling(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	var samp ui.Sampling
	err = json.Unmarshal(body, &samp)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	aH.samplingWriter.WriteSampling(samp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(samp)
}
