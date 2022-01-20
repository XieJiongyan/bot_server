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
	Id             uint
	conn           net.Conn
	ClientOrDevice string
}

// 返回连接构成的 tcp_server, 如果登陆失败，会返回相应的错误
func LoginForConn(conn net.Conn) (TcpServer, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 15))
	// conn.Write([]byte("Connected, start logining in\n"))

	readByte := make([]byte, 1000)
	cnt, err := conn.Read(readByte)
	if err != nil {
		return TcpServer{}, err
	}
	if cnt == 0 {
		err := fmt.Errorf("receive no content for login")
		return TcpServer{}, err
	}

	//登陆，并判断是客户端还是设备端
	isLoginSuccess, id, clientOrDevice := loginAccount(readByte[:cnt])
	if !isLoginSuccess {
		return TcpServer{}, err
	}

	//Sc 增加该 TcpServer
	t := &TcpServer{id, conn, clientOrDevice}
	if clientOrDevice == "client" {
		addClientServer(t, id)
	} else if clientOrDevice == "device" {
		addDeviceServer(t)
		t.conn.Write([]byte("add Device\n"))
	}

	//开启 heartbeat
	go t.heartbeat()
	return *t, err
}

func loginAccount(readByte []byte) (bool, uint, string) {
	clientOrDevice := "unknown"
	type login_content struct {
		Content []string `json:"content"`
	}
	var login_str *login_content = &login_content{}
	login_str.Content = make([]string, 4)
	err := json.Unmarshal(readByte, login_str)
	if len(login_str.Content) < 4 {
		fmt.Println("error: login content shorter than 4 ")
		return false, 0, clientOrDevice
	}
	if err != nil {
		return false, 0, clientOrDevice
	}

	//客户端登陆
	if login_str.Content[0] == "login" && login_str.Content[1] == "client" {
		id_int, err := strconv.Atoi(login_str.Content[2])
		if err != nil {
			return false, 0, clientOrDevice
		}
		id := uint(id_int)

		var password string = login_str.Content[3]

		ok := sc.clientAccount.login(id, password)
		clientOrDevice = "client"
		return ok, id, clientOrDevice
	} else if login_str.Content[0] == "login" && login_str.Content[1] == "device" { //设备登陆
		id_int, err := strconv.Atoi(login_str.Content[2])
		if err != nil {
			return false, 0, clientOrDevice
		}
		id := uint(id_int)

		var password string = login_str.Content[3]

		ok := sc.deviceAccount.login(id, password)
		clientOrDevice = "device"
		return ok, id, clientOrDevice
	}
	fmt.Println("error: unknown login command")
	return false, 0, clientOrDevice
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
	sc.deviceAccount = *loadAccountByJson("data/devices.json")
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
func addDeviceServer(t *TcpServer) bool {
	sc.deviceLock <- true
	defer func() { <-sc.deviceLock }()

	var id uint = t.Id
	_, exist := sc.deviceServer[id]
	if exist {
		fmt.Println("already exist this device")
		return false
	}

	sc.deviceServer[id] = t
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
func removeDeviceServer(id uint) {
	sc.deviceLock <- true
	defer func() { <-sc.deviceLock }()

	if _, isExist := sc.deviceServer[id]; !isExist {
		return
	}
	sc.deviceServer[id].conn.Close()
	sc.deviceAccount.logout(id)
	delete(sc.deviceServer, id)
}

//通过 TcpServer.heartbeat() 确保及时知道网络连接情况
func (t *TcpServer) heartbeat() {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		t.conn.SetDeadline(time.Now().Add(time.Second * 15))
		err := t.Write(NetStruct{Command: "heartbeat"})
		if err != nil {
			fmt.Println(tag, err)
			return
		}
	}
}

// 如果返回错误，应停止该 conn
func (t *TcpServer) Write(n NetStruct) error {
	var b []byte
	b, err := json.Marshal(n)
	if err != nil {
		fmt.Println(tag, "Error marshal: ", err)
		return nil //失败但不停止 conn
	}

	fmt.Println(tag, "Ready to write: ", string(b))
	_, err = t.conn.Write(append(b, '\n'))
	if err != nil {
		if t.ClientOrDevice == "client" {
			removeClientServer(t.Id)
		} else if t.ClientOrDevice == "device" {
			removeDeviceServer(t.Id)
		}
	}
	return err
}

//即使写入失败，也不影响调用的 TcpServer
func WriteToDevice(deviceId uint, n NetStruct) {
	if _, isExist := sc.deviceServer[deviceId]; !isExist {
		return
	}

	sc.deviceServer[deviceId].Write(n)
}

type NetStruct struct {
	Command string   `json:"command"`
	Options []string `json:"options"`
	Extras  string   `json:"extras"`
}

// 阻塞读取，如果返回错误，则说明 Tcp_server 已有错误
func (t *TcpServer) Read() (NetStruct, error) {
	buf := make([]byte, 1000)
	if t.conn == nil {
		err := fmt.Errorf("t.conn is null")
		return NetStruct{}, err
	}
	cnt, err := t.conn.Read(buf)
	if err != nil {
		fmt.Println(tag, "error read: ", err)
		if t.ClientOrDevice == "client" {
			removeClientServer(t.Id)
		} else if t.ClientOrDevice == "device" {
			removeDeviceServer(t.Id)
		}
		return NetStruct{}, err
	}

	is := NetStruct{}
	err = json.Unmarshal(buf[:cnt], &is)
	if err != nil {
		fmt.Println(tag, "unmarshal error: ", err)
		return t.Read()
	}
	return is, err
}
