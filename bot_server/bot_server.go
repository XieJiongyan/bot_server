package botserver

import (
	server "first_server/server"
	"fmt"
	"log"
	"net"
)

func StartTcpServer(port uint) {
	tcp_server, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		log.Panic("having wrong to open tcp")
	}
	fmt.Println("Successfully open tcp_server")

	for {
		coon, err := tcp_server.Accept()
		fmt.Println("Accept tcp server")
		if err != nil {
			fmt.Println("连接出错")
			fmt.Println(err)
		}

		go startCoon(coon)
	}
}

func startCoon(conn net.Conn) {
	_, err := server.LoginForConn(conn)
	if err != nil {
		fmt.Println(err)
	}
}
