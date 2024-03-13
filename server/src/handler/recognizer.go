package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"opencv/src/recognizer"
)

type Recognizer struct {
}

func (h *Recognizer) Recognize(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := recognizer.IsSamePersonById(idParam, imageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if result {
		w.Write([]byte("{result:true}"))
	} else {
		w.Write([]byte("{result:false}"))
	}

}

func (h *Recognizer) AddToRecognize(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = recognizer.RecognizeAndSave(idParam, imageData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Recognizer) RecognizeTwoPhoto(writer http.ResponseWriter, request *http.Request) {
	type bodyStruct struct {
		Example string `json:"example"`
		ToTest  string `json:"to-test"`
	}
	var body bodyStruct
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	example, err := base64.StdEncoding.DecodeString(body.Example)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	toTest, err := base64.StdEncoding.DecodeString(body.ToTest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := recognizer.IsSamePerson(example, toTest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if result {
		writer.Write([]byte("{result:true}"))
	} else {
		writer.Write([]byte("{result:false}"))
	}
}
