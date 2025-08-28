package main

import (
	"fmt"
	"net/http"
)

func (app *application) Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World")
}
