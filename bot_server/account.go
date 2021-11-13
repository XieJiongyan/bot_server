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
	_users map[uint]user
}

func Load_account_by_json(filename string) *account {
	var acc *account = &account{}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(file)
	fmt.Println(byteValue)
	fmt.Println(string(byteValue))
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

	// var tmp_struct map[string]interface{} = make(map[string]interface{})
	// err = json.Unmarshal(byteValue, &tmp_struct)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(tmp_struct)

	// type type1 struct {
	// 	Users float64 `json:"users"`
	// }
	// var t *type1 = &type1{}
	// err = json.Unmarshal(byteValue, t)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(*t)
	return acc
}
