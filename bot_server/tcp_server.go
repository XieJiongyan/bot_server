package botserver

import (
	"fmt"
	"log"
	"net"
)

func coonHandler(c net.Conn) {
	if c == nil {
		log.Panic("connection is nil")
	}
	buf := make([]byte, 1000)
	for {
		cnt, err := c.Read(buf)
		fmt.Println(string(buf))
		if cnt == 0 || err != nil {
			c.Close()
			break
		}
		c.Write([]byte("I see you kotlin"))
	}
}

func Start_tcp_server() {
	tcp_server, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Panic("having wrong to open tcp")
	}
	fmt.Println("Successfully open tcp_server")
	for {
		coon, err := tcp_server.Accept()
		if err != nil {
			fmt.Println("连接出错")
			fmt.Println(err)
		}

		go coonHandler(coon)
	}
}
