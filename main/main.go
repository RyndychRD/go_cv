package main

import (
	"context"
	"fmt"
	"opencv/main/application"
)

func main(){
	app:=application.New()
	err:=app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app:",err)
	}
}
