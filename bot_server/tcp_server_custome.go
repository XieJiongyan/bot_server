package botserver

import "fmt"

func (t *Tcp_server) read() {
	defer Sc.remove_client_server(t.id)

	buf := make([]byte, 1000)
	for {
		cnt, err := t.Conn.Read(buf)
		fmt.Println(string(buf))
		if cnt == 0 {
			continue
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		t.Conn.Write([]byte("I see you kotlin\n"))
	}
}

func (t *Tcp_server) write() {

}
