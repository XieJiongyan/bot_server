package botserver

import (
	"fmt"
	"time"
)

func (t *Tcp_server) customer_start() {
	go t.read()
	go t.heartbeat()
}

func (t *Tcp_server) read() {
	defer Sc.remove_client_server(t.id)

	buf := make([]byte, 1000)
	for {
		cnt, err := t.Conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf))
		if cnt == 0 {
			continue
		}
		t.Conn.Write([]byte("I see you kotlin\n"))
	}
}

func (t *Tcp_server) heartbeat() {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		t.Conn.Write([]byte("heartbeat\n"))
	}
}
