package botserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type server_center struct {
	client_account account
	device_account account

	client_server map[uint]*Tcp_server
	device_server map[uint]*Tcp_server

	client_max_cnt uint
	device_max_cnt uint

	client_lock chan bool
	device_lock chan bool
}

var sc *server_center

func (sc *server_center) add_client_server(t *Tcp_server, id uint, password string) bool {
	sc.client_lock <- true
	defer func() { <-sc.client_lock }()

	t.id = id
	_, exist := sc.client_server[id]
	if exist {
		fmt.Println("already exist unit")
		return false
	}

	if ok := sc.client_account.Login(id, password); !ok {
		return false
	}

	sc.client_server[id] = t
	return true
}

func (sc *server_center) remove_client_server(id uint) {
	sc.client_lock <- true
	defer func() { <-sc.client_lock }()
	sc.client_server[id]._coon.Close()
	delete(sc.client_server, id)
}

type Tcp_server struct {
	id    uint
	_coon net.Conn
}

const _OUTPUT_SIZE uint = 10

func start_server(coon net.Conn) {
	var t *Tcp_server = &Tcp_server{}
	if t._coon == nil {
		log.Panic("connection is nil")
	}
	defer t._coon.Close()

	t._coon.SetReadDeadline(time.Now().Add(time.Second * 180))
	read_byte := make([]byte, 1000)
	cnt, err := t._coon.Read(read_byte)
	fmt.Println(string(read_byte))
	if cnt == 0 || err != nil {
		return
	}

	type login_content struct {
		content []string
	}
	var login_str *login_content = &login_content{}
	login_str.content = make([]string, 4)
	err = json.Unmarshal(read_byte, login_str)
	if err != nil || len(login_str.content) < 4 {
		fmt.Println(err)
	}

	// 根据 login_str 中的内容判断动作，如果全部匹配不到，会关闭连接并关闭函数
	// 这里是建联过程，如果失败，将会断开连接
	if login_str.content[0] == "login" {
		if login_str.content[1] == "client" {
			id_int, err := strconv.Atoi(login_str.content[2])
			if err != nil {
				fmt.Println(err)
			}
			t.id = uint(id_int)

			var password string = login_str.content[3]
			sc.add_client_server(t, t.id, password)
			defer sc.remove_client_server(t.id)

		}
	}

}
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
		c.Write([]byte("I see you kotlin\n"))
	}
}

func Start_tcp_server(port uint) {
	tcp_server, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		log.Panic("having wrong to open tcp")
	}
	fmt.Println("Successfully open tcp_server")

	sc = &server_center{}
	sc.client_account = *Load_account_by_json("data/clients.json")
	sc.device_account = *account_construct()
	for {
		coon, err := tcp_server.Accept()
		fmt.Println("Accept tcp server")
		if err != nil {
			fmt.Println("连接出错")
			fmt.Println(err)
		}

		go coonHandler(coon)
	}
}
