package botserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type user struct {
	Id         uint   `json:"id"`
	Password   string `json:"password"`
	Is_login   bool   `json:"is_login"`
	login_lock chan bool
}

type account struct {
	_users     map[uint]user
	write_lock chan bool
}

func account_construct() *account {
	var acc *account = &account{}
	acc._users = make(map[uint]user)
	acc.write_lock = make(chan bool, 1)
	return acc
}
func Load_account_by_json(filename string) *account {
	var acc *account = &account{}
	acc._users = make(map[uint]user)
	acc.write_lock = make(chan bool, 1)

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	type tmp_acc_type struct {
		Users []user `json:"users"`
	}
	var tmp_acc *tmp_acc_type = &tmp_acc_type{}
	err = json.Unmarshal(byteValue, tmp_acc)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tmp_acc)

	for _, u := range tmp_acc.Users {
		acc.add_user(u)
	}
	return acc
}

// 并发安全地为账号系统增添用户
// 只有在 account 中不含有 u 的 id 的 user 时， 才会增添账号
// 增添的账号确保处于非登陆的状态
// 增添的账号处于非锁定状态
func (a *account) add_user(u user) {
	var id uint = u.Id
	a.write_lock <- true
	_, exists := a._users[id]
	if exists {
		<-a.write_lock
		return
	}
	u_copy := u
	u_copy.Is_login = false
	u_copy.login_lock = make(chan bool, 1)
	a._users[id] = u_copy
	<-a.write_lock
}

func (a *account) Login(id uint, password string) bool {
	u, exist := a._users[id]
	if !exist {
		return false
	}
	u.login_lock <- true
	if u.Is_login || password != u.Password {
		fmt.Println("Already login or false Password")
		<-u.login_lock
		return false
	}
	u.Is_login = true
	<-u.login_lock
	return true
}

func (a *account) Logout(id uint) {
	u, exist := a._users[id]
	if !exist {
		return
	}
	u.login_lock <- true
	u.Is_login = false
	<-u.login_lock
}
