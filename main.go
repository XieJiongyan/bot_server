package main

import (
	botserver "first_server/bot_server"
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
	fmt.Println(string(botserver.Get_devices_byte()))
	io.WriteString(w, string(botserver.Get_devices_byte()))
}

func start_http_server() {
	fmt.Println(string(botserver.Get_devices_byte()))

	handler := my_handler{}
	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func main() {
	go start_http_server()
	go botserver.Start_tcp_server()
	go botserver.Load_account_by_json("data/clients.json")
	select {}
}
