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
	modelsDir       = filepath.Join("recognizer", "models")
	imagesDir       = filepath.Join("recognizer", "images")
	isTestThreshold = false
)

func InitEnv() {
	modelPath, mExists := os.LookupEnv("MODEL_FULL_PATH")
	if mExists && modelPath != "" {
		fmt.Println("model path set")
		modelsDir = modelPath
	}
	storagePath, sExists := os.LookupEnv("STORAGE_FULL_PATH")
	if sExists && storagePath != "" {
		fmt.Println("image path set")
		imagesDir = storagePath
	}
	isTest, tExists := os.LookupEnv("IS_TEST_THRESHOLD")
	if tExists && isTest == "1" {
		fmt.Println("is test threshold set")
		isTestThreshold = true
	}
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
		Err = &MyError{
			custom: "can't init face recognizer:",
			origin: err,
		}
		return
	}
	defer rec.Close()

	directory, Err := checkTrainingDataExist(id)
	if Err != nil {
		return
	}
	samples, cats, Err := train(directory, rec)
	if Err != nil {
		return
	}
	// get face from in image
	f, Err := recognizeOneFace(rec, image)
	if Err != nil {
		return
	}
	rec.SetSamples(samples, cats)
	catID := rec.ClassifyThreshold(f.Descriptor, float32(0.2))
	testThresholdIfEnabled(rec, f)
	return catID == 200, nil
}

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

func train(directory string, rec *face.Recognizer) (samples []face.Descriptor, cats []int32, Err error) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
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
				cats = append(cats, 200)
			}
		}
		return nil
	})
	if err != nil {
		Err = &MyError{custom: "error while collecting training data:", origin: err}
		return
	}
	return samples, cats, err
}

func checkTrainingDataExist(id string) (directory string, Err error) {
	directory = filepath.Join(imagesDir, id)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		Err = errors.New("no training data found")
		return
	}
	return
}

func testThresholdIfEnabled(rec *face.Recognizer, f face.Face) {
	if isTestThreshold {
		for i := 1; i <= 10; i++ {
			thr := float32(i) / float32(10)
			catID := rec.ClassifyThreshold(f.Descriptor, thr)
			fmt.Printf("cat by classificator %v with threshold %f \n", catID, thr)
		}
	}
}
