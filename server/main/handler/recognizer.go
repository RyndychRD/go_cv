package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"opencv/main/recognizer"
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

	result, err := recognizer.IsSamePerson(idParam, imageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(result)
	if result == true {
		w.Write([]byte("{status:true}"))
	} else {
		w.Write([]byte("{status:false}"))
	}
	w.WriteHeader(http.StatusOK)

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