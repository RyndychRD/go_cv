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