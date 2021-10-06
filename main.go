package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type my_handler struct{}

type Ret struct {
	Message int
	Str     string
}

func (my_handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The method is", r.Method)
	ret := Ret{
		Message: 100,
		Str:     "str",
	}
	ret_json, _ := json.Marshal(ret)
	fmt.Println(string(ret_json))
	io.WriteString(w, string(ret_json))
}

func main() {
	handler := my_handler{}
	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
