package handler

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"mime/multipart"
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

func getExampleAndToTestByteArray(request *http.Request) (example []byte, toTest []byte, err error) {
	err = request.ParseMultipartForm(10 << 20) // Максимальный размер 32MB
	if err != nil {
		err = errors.New(fmt.Sprintf("error parsing form, error: %+v", err))
		return
	}

	exampleFile, _, err := request.FormFile("example")
	if err != nil {
		err = errors.New(fmt.Sprintf("error parsing example file, error: %+v", err))
		return
	}
	defer func(exampleFile multipart.File) {
		err := exampleFile.Close()
		if err != nil {
			log.Printf("error closing file, error: %+v", err)
		}
	}(exampleFile)

	toTestFile, _, err := request.FormFile("to-test")
	if err != nil {
		err = errors.New(fmt.Sprintf("error parsing to-test file, error: %+v", err))
		return
	}
	defer func(toTestFile multipart.File) {
		err := toTestFile.Close()
		if err != nil {
			log.Printf("error closing file, error: %+v", err)
		}
	}(toTestFile)

	example, err = io.ReadAll(exampleFile)
	if err != nil {
		err = errors.New(fmt.Sprintf("error parsing example file to byte array, error: %+v", err))
		return
	}

	toTest, err = io.ReadAll(toTestFile)
	if err != nil {
		err = errors.New(fmt.Sprintf("error parsing to-test file to byte array, error: %+v", err))
		return
	}
	return
}

func (h *Recognizer) RecognizeTwoPhoto(writer http.ResponseWriter, request *http.Request) {
	example, toTest, err := getExampleAndToTestByteArray(request)
	if err != nil {
		badRequest(writer, fmt.Sprintf("error reading form data, error: %+v", err))
	}
	if example, err = convertion.ToJpeg(example); err != nil {
		badRequest(writer, fmt.Sprintf("can't convert example image to jpeg, error: %+v", err))
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
