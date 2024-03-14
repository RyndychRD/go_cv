#### репозиторий
https://github.com/Kagami/go-face

### репозиторий с тестовыми данными
https://github.com/Kagami/go-face-testdata/tree/master

### тутор с картинками
https://tutorialedge.net/golang/go-face-recognition-tutorial-part-one/


### Команда powershell для перегона в base64

```cmd
powershell -command "& {[System.Convert]::ToBase64String([System.IO.File]::ReadAllBytes('D:\projects\go_opencv\server\images\nurbek_doc.jpg'))}" > nurbek_doc.txt
```

### Различные решения на Go для распознавания и идентификации лиц

#### Только распознавание точек на лицах
1. https://github.com/modanesh/GoFaceRec
2. https://github.com/esimov/pigo/tree/master
3. https://github.com/hybridgroup/gocv

#### Идентификация лиц
1. https://github.com/exadel-inc/CompreFace
   2. Запустить так и не смог, ресурсы жрет как не в себя. Запускал в wsl с 4 процами и 4 гигами оперативки
2. https://deepstack.readthedocs.io/en/latest/index.html
   3. Некоторое стороннее решение, не пробовал, только нашел
   4. https://gocv.io/writing-code/more-examples/
