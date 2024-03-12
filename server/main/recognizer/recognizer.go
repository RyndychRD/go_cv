package recognizer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	modelsDir = filepath.Join("models")
	imagesDir = filepath.Join("images")
)

func RecognizeMain() {
	// Init the recognizer.
	rec, err := face.NewRecognizer(modelsDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	// Test image with 10 faces.
	testImagePristin := filepath.Join(imagesDir, "pristin.jpg")
	// Recognize faces on that image.
	faces, err := rec.RecognizeFile(testImagePristin)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	if len(faces) != 10 {
		log.Fatalf("Wrong number of faces")
	}

	testImage1 := filepath.Join(imagesDir, "azamat.jpg")
	testImage2 := filepath.Join(imagesDir, "azamat2.jpg")
	face1, err := rec.RecognizeSingleFile(testImage1)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}
	face2, err := rec.RecognizeSingleFile(testImage2)
	if err != nil {
		log.Fatalf("Can't recognize: %v", err)
	}

	var samples []face.Descriptor
	var cats []int32
	for i, f := range faces {
		samples = append(samples, f.Descriptor)
		// Each face is unique on that image so goes to its own category.
		cats = append(cats, int32(i))
	}
	samples = append(samples, face1.Descriptor)
	samples = append(samples, face2.Descriptor)
	//fmt.Println("Face3 start")
	//samples = append(samples, face3.Descriptor)

	cats = append(cats, 10, 10)

	folderPath := filepath.Join(imagesDir, "me")
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Если это файл (а не папка), печатаем его путь
		if !info.IsDir() {
			faceT, err := rec.RecognizeSingleFile(path)
			if err != nil {
				log.Fatalf("Can't recognize: %v", err)
			}
			if faceT != nil {
				samples = append(samples, faceT.Descriptor)
				cats = append(cats, 11)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Ошибка при обходе файлов:", err)
	}

	// Name the categories, i.e. people on the image.
	labels := []string{
		"Sungyeon", "Yehana", "Roa", "Eunwoo", "Xiyeon",
		"Kyulkyung", "Nayoung", "Rena", "Kyla", "Yuha", "Azamat1", "Roman",
	}

	rec.SetSamples(samples, cats)
	fmt.Println("Samples set")

	recognize(rec, labels, "me_test.jpg")
	recognize(rec, labels, "me_test_front_box.jpg")
	recognize(rec, labels, "me_test_front_hair.jpg")
	recognize(rec, labels, "me_test_front_hand.jpg")
	recognize(rec, labels, "me_test_up_clear.jpg")
	recognize(rec, labels, "me_test_up_hair.jpg")
}

func recognize(rec *face.Recognizer, labels []string, path string) {
	testImageNayoung := filepath.Join(imagesDir, path)
	f, err := rec.RecognizeSingleFile(testImageNayoung)
	if err != nil {
		fmt.Printf("Can't recognize  %s: %v", path, err)
		return
	}
	if f == nil {
		fmt.Printf("Not a single face on the image %s", path)
		return
	}
	catID := rec.Classify(f.Descriptor)
	if catID < 0 {
		fmt.Printf("Can't classify %s", path)
		return
	}
	fmt.Println(path, " ", labels[catID])
}
