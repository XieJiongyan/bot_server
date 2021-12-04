// tcp_server 将与 account 一起实现登陆和账号管理功能
// 支持 Login(conn) TcpServer, TcpServer.Write(), TcpServer.Read() 三个功能
package server

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type TcpServer struct {
	Id   uint
	conn net.Conn
}

// 返回连接构成的 tcp_server, 如果登陆失败，会返回相应的错误
func LoginForConn(conn net.Conn) (TcpServer, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 15))
	conn.Write([]byte("Connected, start logining in\n"))

	readByte := make([]byte, 1000)
	cnt, err := conn.Read(readByte)
	if err != nil {
		return TcpServer{}, err
	}
	if cnt == 0 {
		err := fmt.Errorf("receive no content for login")
		return TcpServer{}, err
	}

	//
	isLoginSuccess, id, err := loginAccount(readByte[:cnt])
	if !isLoginSuccess {
		return TcpServer{}, err
	}

	//Sc 增加该 TcpServer
	t := &TcpServer{id, conn}
	addClientServer(t, id)

	//开启 heartbeat
	go t.heartbeat()
	return *t, err
}

func loginAccount(readByte []byte) (bool, uint, error) {
	type login_content struct {
		Content []string `json:"content"`
	}
	var login_str *login_content = &login_content{}
	login_str.Content = make([]string, 4)
	err := json.Unmarshal(readByte, login_str)
	if len(login_str.Content) < 4 {
		err = fmt.Errorf("login content shorter than 4 ")
		return false, 0, err
	}
	if err != nil {
		return false, 0, err
	}

	if login_str.Content[0] == "login" && login_str.Content[1] == "client" {
		id_int, err := strconv.Atoi(login_str.Content[2])
		if err != nil {
			return false, 0, err
		}
		id := uint(id_int)

		var password string = login_str.Content[3]

		ok := sc.clientAccount.login(id, password)
		return ok, id, err
	}
	return false, 0, err
}

type server_center struct {
	clientAccount account
	deviceAccount account

	clientServer map[uint]*TcpServer
	deviceServer map[uint]*TcpServer

	clientMaxCnt uint
	deviceMaxCnt uint

	clientLock chan bool
	deviceLock chan bool
}

/**
单例模式
*/
var sc *server_center

func init() {
	sc = &server_center{}
	sc.clientAccount = *loadAccountByJson("data/clients.json")
	sc.deviceAccount = *account_construct()
	sc.clientLock = make(chan bool, 1)
	sc.deviceLock = make(chan bool, 1)
	sc.clientServer = make(map[uint]*TcpServer)
	sc.deviceServer = make(map[uint]*TcpServer)
}

func addClientServer(t *TcpServer, id uint) bool {
	sc.clientLock <- true
	defer func() { <-sc.clientLock }()

	t.Id = id
	_, exist := sc.clientServer[id]
	if exist {
		fmt.Println("already exist unit")
		return false
	}

	sc.clientServer[id] = t
	return true
}

func removeClientServer(id uint) {
	sc.clientLock <- true
	defer func() { <-sc.clientLock }()
	if _, isExist := sc.clientServer[id]; !isExist {
		return
	}
	sc.clientServer[id].conn.Close()
	sc.clientAccount.logout(id)
	delete(sc.clientServer, id)
}

//通过 TcpServer.heartbeat() 确保及时知道网络连接情况
func (t *TcpServer) heartbeat() {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		t.conn.SetDeadline(time.Now().Add(time.Second * 15))
		err := t.Write([]byte("heartbeat\n"))
		if err != nil {
			fmt.Println(tag, err)
			return
		}
	}
}

func (t *TcpServer) Write(b []byte) error {
	_, err := t.conn.Write(b)
	if err != nil {
		removeClientServer(t.Id)
	}
	return err
}

type InputStruct struct {
	Command string                 `json:"command"`
	Options []string               `json:"options"`
	Extras  map[string]interface{} `json:"extras"`
}

// 阻塞读取，如果返回错误，则说明 Tcp_server 已有错误
func (t *TcpServer) Read() (InputStruct, error) {
	buf := make([]byte, 1000)
	_, err := t.conn.Read(buf)
	if err != nil {
		removeClientServer(t.Id)
		return InputStruct{}, err
	}

	is := InputStruct{}
	err = json.Unmarshal(buf, &is)
	if err != nil {
		fmt.Println(tag, "unmarshal error: ", err)
		return t.Read()
	}
	return is, err
}
