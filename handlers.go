package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sshetty10/go-seed-db/model"
)

func (a *API) ListTrainers(w http.ResponseWriter, r *http.Request) {

	trainers, err := a.ListDBTrainers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, trainers)
}

func (a *API) GetTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	t, err := a.GetDBTrainerByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *API) DeleteTrainer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	err := a.DeleteDBTrainer(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, "Trainer Deleted")
}

func (a *API) CreateTrainer(w http.ResponseWriter, r *http.Request) {
	var t model.Trainer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	err := a.CreateDBTrainer(&t)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, t)
}

func (a *API) SayCheese(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]string)
	m["msg"] = "Say Cheese"
	respondWithJSON(w, http.StatusOK, m)
}

func (a *API) GetTrainerFirstName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]

	t, err := a.GetDBTrainerByName(name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

// Response structure
type Response struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	r := &Response{Errors: []string{message}}
	response, _ := json.Marshal(r)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response := &Response{Data: payload}
	data, err := json.Marshal(response)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
