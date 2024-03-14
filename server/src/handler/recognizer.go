package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"opencv/src/convertion"
	"opencv/src/recognizer"
)

type Recognizer struct {
}

func (h *Recognizer) Recognize(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	imageData, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := recognizer.IsSamePersonById(idParam, imageData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if result {
		if _, err := writer.Write([]byte("{result:true}")); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := writer.Write([]byte("{result:false}")); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func (h *Recognizer) AddToRecognize(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	imageData, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = recognizer.RecognizeAndSave(idParam, imageData); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (h *Recognizer) RecognizeTwoPhoto(writer http.ResponseWriter, request *http.Request) {
	type bodyStruct struct {
		Example string `json:"example"`
		ToTest  string `json:"to-test"`
	}
	var body bodyStruct
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	example, err := base64.StdEncoding.DecodeString(body.Example)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	if example, err = convertion.ToJpeg(example); err != nil {
		http.Error(writer, fmt.Sprintf("bad example picture: %+v", err), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	toTest, err := base64.StdEncoding.DecodeString(body.ToTest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	if toTest, err = convertion.ToJpeg(toTest); err != nil {
		http.Error(writer, fmt.Sprintf("bad to-test picture: %+v", err), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	result, err := recognizer.IsSamePerson(example, toTest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println(err.Error())
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if result {
		if _, err := writer.Write([]byte("{result:true}")); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	} else {
		if _, err := writer.Write([]byte("{result:false}")); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	}
}
