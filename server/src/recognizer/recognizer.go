package recognizer

import (
	"errors"
	"fmt"
	"github.com/Kagami/go-face"
	"github.com/google/uuid"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	modelsDir                       = filepath.Join("recognizer", "models")
	imagesDir                       = filepath.Join("recognizer", "images")
	isTestThreshold                 = false
	catForIdentification            = 200
	ClassificationThreshold float32 = 0.3
)

func InitEnv() {
	modelPath, mExists := os.LookupEnv("MODEL_FULL_PATH")
	if mExists && modelPath != "" {
		modelsDir = modelPath
	} else {
		log.Println("bad MODEL_FULL_PATH, using default ", modelsDir)
	}
	storagePath, sExists := os.LookupEnv("STORAGE_FULL_PATH")
	if sExists && storagePath != "" {
		imagesDir = storagePath
	} else {
		log.Println("bad STORAGE_FULL_PATH, using default ", imagesDir)
	}
	isTest, tExists := os.LookupEnv("IS_TEST_THRESHOLD")
	if tExists && isTest == "1" {
		isTestThreshold = true
	} else {
		log.Println("bad IS_TEST_THRESHOLD, using default ", isTestThreshold)
	}
	threshold, thresExists := os.LookupEnv("THRESHOLD_VALUE")
	if thresExists {
		if s, err := strconv.ParseFloat(threshold, 32); err == nil {
			ClassificationThreshold = float32(s)
		} else {
			log.Println("bad THRESHOLD_VALUE, using default ", ClassificationThreshold)
		}
	} else {
		log.Println("bad THRESHOLD_VALUE, using default ", ClassificationThreshold)
	}
}

func RecognizeAndSave(id string, image []byte) error {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("can't init face recognizer: %v", err)
	}
	defer rec.Close()
	if _, err = recognizeOneFace(rec, image, false); err != nil {
		return err
	}
	if err = saveFile(id, image); err != nil {
		return err
	}
	return nil
}

func IsSamePersonById(id string, image []byte) (result bool, Err error) {
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
	f, Err := recognizeOneFace(rec, image, false)
	if Err != nil {
		return
	}
	rec.SetSamples(samples, cats)
	catID := rec.ClassifyThreshold(f.Descriptor, ClassificationThreshold)
	testThresholdIfEnabled(rec, f)
	return catID == catForIdentification, nil
}

func IsSamePerson(example []byte, toTest []byte) (result bool, Err error) {
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		Err = &MyError{
			custom: "can't init face recognizer:",
			origin: err,
		}
		return
	}
	defer rec.Close()

	exampleFace, err := recognizeOneFace(rec, example, true)
	if err != nil {
		Err = &MyError{
			custom: "bad example face",
			origin: err,
		}
		return
	}
	toTestFace, err := recognizeOneFace(rec, toTest, false)
	if err != nil {
		Err = &MyError{
			custom: "bad to-test face",
			origin: err,
		}
		return
	}
	samples := []face.Descriptor{exampleFace.Descriptor}
	cats := []int32{int32(catForIdentification)}
	rec.SetSamples(samples, cats)
	catID := rec.ClassifyThreshold(toTestFace.Descriptor, ClassificationThreshold)
	testThresholdIfEnabled(rec, toTestFace)
	return catID == catForIdentification, nil
}

type MyError struct {
	custom string
	origin error
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s, origin: %s",
		e.custom, e.origin)
}

func recognizeOneFace(rec *face.Recognizer, image []byte, isIgnoreSeveralFaces bool) (face face.Face, Err error) {
	faces, err := rec.Recognize(image)
	if err != nil {
		Err = &MyError{custom: "can't recognize face", origin: err}
		return
	}
	if len(faces) == 0 {
		Err = errors.New("found no face")
		return
	}
	if isIgnoreSeveralFaces == false && len(faces) > 1 {
		Err = errors.New("found more than 1 face on single photo")
		return
	}
	face = faces[0]
	return
}

func saveFile(id string, image []byte) error {
	directory := filepath.Join(imagesDir, id)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			return &MyError{custom: "can't create destination folder:", origin: err}
		}
	}
	if err := os.WriteFile(filepath.Join(imagesDir, id, uuid.New().String()+".jpg"), image, 0777); err != nil {
		return &MyError{custom: "can't save new image to destination:", origin: err}
	}
	return nil
}

func train(directory string, rec *face.Recognizer) (samples []face.Descriptor, cats []int32, Err error) {
	catForIdentificationInt32 := int32(catForIdentification)
	var paths []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			lastIndex := strings.LastIndex(path, ".")
			if lastIndex != -1 && path[lastIndex+1:] == "jpg" {
				paths = append(paths, path)
			}
		}
		return nil
	})
	if err != nil {
		Err = &MyError{custom: "error while collecting training data:", origin: err}
		return
	}
	if len(paths) == 0 {
		Err = errors.New("no photos found in training data")
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(len(paths))
	for i := 0; i < len(paths); i++ {
		go func() {
			defer wg.Done()

			var path string
			path, paths = paths[len(paths)-1], paths[:len(paths)-1]
			faceT, err := rec.RecognizeSingleFile(path)
			if err != nil {
				log.Println(fmt.Sprintf("error while recognizing %s, origin: %v", path, err))
			}
			if faceT != nil {
				mu.Lock()
				samples = append(samples, faceT.Descriptor)
				cats = append(cats, catForIdentificationInt32)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
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
		for i := 5; i <= 100; i = i + 1 {
			thr := float32(i) / float32(100)
			catID := rec.ClassifyThreshold(f.Descriptor, thr)
			fmt.Printf("cat by classificator %v with threshold %f \n", catID, thr)
			if catID == catForIdentification {
				return
			}
		}
	}
}
