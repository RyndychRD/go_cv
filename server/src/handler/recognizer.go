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
		badRequest(writer, fmt.Sprintf("can't read received binary, error: %+v", err))
		return
	}
	if imageData, err = convertion.ToJpeg(imageData); err != nil {
		badRequest(writer, fmt.Sprintf("can't convert image to jpeg, error: %+v", err))
		return
	}

	result, thr, err := recognizer.IsSamePersonById(idParam, imageData)
	if err != nil {
		badRequest(writer, fmt.Sprintf("can't recognize person, error: %+v", err))
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var response string
	if result {
		response = fmt.Sprintf("{result:true, thr:%f}", thr)
	} else {
		response = "{result:false}"
	}
	if _, err := writer.Write([]byte(response)); err != nil {
		serverError(writer, fmt.Sprintf("error writing results, error: %+v", err))
		return
	}

}

func (h *Recognizer) AddToRecognize(writer http.ResponseWriter, request *http.Request) {
	idParam := chi.URLParam(request, "id")
	imageData, err := io.ReadAll(request.Body)
	if err != nil {
		badRequest(writer, fmt.Sprintf("can't read received binary, error: %+v", err))
		return
	}
	if imageData, err = convertion.ToJpeg(imageData); err != nil {
		badRequest(writer, fmt.Sprintf("can't convert image to jpeg, error: %+v", err))
		return
	}

	if err = recognizer.RecognizeAndSave(idParam, imageData); err != nil {
		badRequest(writer, fmt.Sprintf("can't add photo to train data, error: %+v", err))
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
		badRequest(writer, fmt.Sprintf("can't decode json, error: %+v", err))
		return
	}

	example, err := base64.StdEncoding.DecodeString(body.Example)
	if err != nil {
		badRequest(writer, fmt.Sprintf("can't decode example image from base64, error: %+v", err))
		return
	}
	if example, err = convertion.ToJpeg(example); err != nil {
		badRequest(writer, fmt.Sprintf("can't convert example image to jpeg, error: %+v", err))
		return
	}

	toTest, err := base64.StdEncoding.DecodeString(body.ToTest)
	if err != nil {
		badRequest(writer, fmt.Sprintf("can't decode to-test image from base64, error: %+v", err))
		return
	}
	if toTest, err = convertion.ToJpeg(toTest); err != nil {
		badRequest(writer, fmt.Sprintf("can't convert to-test image to jpeg, error: %+v", err))
		return
	}

	result, thr, err := recognizer.IsSamePerson(example, toTest)
	if err != nil {
		badRequest(writer, fmt.Sprintf("error comparing faces, error: %+v", err))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	var response string
	if result {
		response = fmt.Sprintf("{result:true, thr:%f}", thr)
	} else {
		response = "{result:false}"
	}
	if _, err := writer.Write([]byte(response)); err != nil {
		serverError(writer, fmt.Sprintf("error writing results, error: %+v", err))
		return
	}
}

func badRequest(writer http.ResponseWriter, response string) {
	http.Error(writer, response, http.StatusBadRequest)
	log.Println(response)
}
func serverError(writer http.ResponseWriter, response string) {
	http.Error(writer, response, http.StatusInternalServerError)
	log.Println(response)
}
