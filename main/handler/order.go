package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o *Order) Create(w http.ResponseWriter, r *http.Request){
	fmt.Println("Create new order")
}

func (o *Order) List(w http.ResponseWriter, r *http.Request){
	fmt.Println("Get all orders")
}
func (o *Order) GetByID(w http.ResponseWriter, r *http.Request){
	fmt.Println("Get one order")
}
func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request){
	fmt.Println("Update one order")
}
func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request){
	fmt.Println("Delete one order")
}