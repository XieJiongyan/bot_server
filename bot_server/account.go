package botserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type user struct {
	Id       uint   `json:"id"`
	Password string `json:"password"`
	Is_login bool   `json:"is_login"`
}

type account struct {
	_users     map[uint]user
	write_lock chan bool
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
func (a *account) add_user(u user) {
	var id uint = u.Id
	a.write_lock <- true
	_, exists := a._users[id]
	if exists {
		<-a.write_lock
		return
	}
	a._users[id] = u
	<-a.write_lock
}
