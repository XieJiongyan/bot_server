package botserver

import (
	"fmt"
	"io/ioutil"
	"os"
)

type alarm struct {
	corn string
	text string
}
type device_info struct {
	user_name string
	password  string
	alarms    []alarm
}

type Server_data struct {
	devices      []device_info
	devices_byte []byte
}

var server_data Server_data

func init() {
	json_file, err := os.Open("bot_data.json")
	if err != nil {
		fmt.Println(err)
	}
	defer json_file.Close()

	server_data.devices_byte, err = ioutil.ReadAll(json_file)
	if err != nil {
		fmt.Print(err)
	}
}

func Get_devices_byte() []byte {
	return server_data.devices_byte
}
