package recognizer

import (
	"errors"
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
)

var (
	modelsDir = filepath.Join("recognizer", "models")
	imagesDir = filepath.Join("recognizer", "images")
)

type MyError struct {
	custom string
	origin error
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s, origin: %s",
		e.custom, e.origin)
}

func recognizeOneFace(rec *face.Recognizer, image []byte) (face face.Face, Err error) {
	faces, err := rec.Recognize(image)
	if err != nil {
		Err = &MyError{custom: "can't recognize face", origin: err}
		return
	}
	if len(faces) == 0 {
		Err = errors.New("found no face ")
		return
	}
	if len(faces) > 1 {
		Err = errors.New("found more than 1 face ")
		return
	}
	face = faces[0]
	return
}

func saveFile(id string, image []byte) error {
	directory := filepath.Join(imagesDir, id)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.Mkdir(directory, 0755)
		if err != nil {
			return &MyError{custom: "can't create destination folder:", origin: err}
		}
	}
	os.WriteFile(filepath.Join(imagesDir, id, uuid.New().String()+".jpg"), image, 0777)
	return nil
}

func RecognizeAndSave(id string, image []byte) error {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("can't init face recognizer: %v", err)
	}
	defer rec.Close()
	if _, err = recognizeOneFace(rec, image); err != nil {
		return err
	}
	if err = saveFile(id, image); err != nil {
		return err
	}
	return nil
}

func IsSamePerson(id string, image []byte) (result bool, Err error) {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("can't init face recognizer: %v", err)
	}
	defer rec.Close()

	directory := filepath.Join(imagesDir, id)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		Err = errors.New("no training data found")
		return
	}

	var samples []face.Descriptor
	var cats []int32
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			faceT, err := rec.RecognizeSingleFile(path)
			if err != nil {
				return err
			}
			if faceT != nil {
				samples = append(samples, faceT.Descriptor)
				cats = append(cats, 0)
			}
		}
		return nil
	})

	if err != nil {
		Err = &MyError{custom: "error while collecting training data:", origin: err}
		return
	}

	f, Err := recognizeOneFace(rec, image)
	if Err != nil {
		return
	}
	catID := rec.Classify(f.Descriptor)
	fmt.Println(catID)
	return catID == 0, nil

}
