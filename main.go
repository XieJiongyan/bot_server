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
	botserver.Start_tcp_server()
	select {}
}

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// )

// // TCP Server端测试
// // 处理函数
// func process(conn net.Conn) {
// 	defer conn.Close() // 关闭连接
// 	for {
// 		reader := bufio.NewReader(conn)
// 		var buf [128]byte
// 		n, err := reader.Read(buf[:]) // 读取数据
// 		if err != nil {
// 			fmt.Println("read from client failed, err: ", err)
// 			break
// 		}
// 		recvStr := string(buf[:n])
// 		fmt.Println("收到Client端发来的数据：", recvStr)
// 		conn.Write([]byte(recvStr)) // 发送数据
// 	}
// }

// func main() {
// 	listen, err := net.Listen("tcp", ":8100")
// 	fmt.Println("aaaaaaa")
// 	if err != nil {
// 		fmt.Println("Listen() failed, err: ", err)
// 		return
// 	}
// 	for {
// 		conn, err := listen.Accept() // 监听客户端的连接请求
// 		if err != nil {
// 			fmt.Println("Accept() failed, err: ", err)
// 			continue
// 		}
// 		go process(conn) // 启动一个goroutine来处理客户端的连接请求
// 	}
// }
