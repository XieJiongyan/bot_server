package server

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )

// const tag = "SERVER"

// func (t *botserver.Tcp_server) customer_start() {
// 	go t.read()
// 	go t.heartbeat()
// }

// type InputJson struct {
// 	Command string      `json:"command"`
// 	Option  []string    `json:"Option"`
// 	Extras  interface{} `json:"extras"`
// }

// func (t *botserver.Tcp_server) read() {
// 	defer Sc.remove_client_server(t.id)

// 	buf := make([]byte, 1000)
// 	for {
// 		cnt, err := t.Conn.Read(buf)
// 		if err != nil {
// 			fmt.Println(tag, err)
// 			return
// 		}
// 		fmt.Println(tag, string(buf))
// 		if cnt == 0 {
// 			continue
// 		}
// 		var input InputJson
// 		err = json.Unmarshal(buf, &input)
// 		if err != nil {
// 			fmt.Println(tag, err)
// 		}
// 		if input.Command == "clock" {
// 			go clock.DealClock(t.Conn, input)
// 		}
// 	}
// }

// func (t *botserver.Tcp_server) heartbeat() {
// 	for {
// 		time.Sleep(time.Duration(5) * time.Second)
// 		t.Conn.Write([]byte("heartbeat\n"))
// 	}
// }
