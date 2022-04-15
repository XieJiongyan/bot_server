//
package clock

import (
	"encoding/json"
	"first_server/server"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const tag = "clock "

type clock struct {
	Cron     string `json:"cron"`
	Text     string `json:"text"`
	IsActive bool   `json:"is_active"`
}

type device struct {
	Name   string  `json:"name"`
	Clocks []clock `json:"clocks"`
}

type client struct {
	Devices map[string]device `json:"devices"`
}

type clockData struct {
	Clients map[string]client `json:"clients"`
}

var cd *clockData

const clockFile = "data/clock.json"

func init() {
	cd = &clockData{}

	file, err := os.Open(clockFile)
	if err != nil {
		fmt.Println(tag, "error open file: ", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(tag, "error read file: ", err)
	}
	err = json.Unmarshal(bytes, cd)
	if err != nil {
		fmt.Println(tag, "error unmarshal: ", err)
	}
}

//处理客户端请求
//get: 得到客户端的所有闹钟信息
//post: 修改某一闹钟
func DealClockForClient(is server.NetStruct, t server.TcpServer) error {
	if len(is.Options) >= 1 && is.Options[0] == "get" {
		writeByte, err := json.Marshal(cd.Clients[fmt.Sprint(t.Id)])
		if err != nil {
			fmt.Println(tag, "error marshel: ", err)
			return nil
		}
		writeByte = append(writeByte, byte('\n'))

		nets := server.NetStruct{
			Command: "clock",
			Options: []string{"total"},
			Extras:  string(writeByte),
		}
		err = t.Write(nets)
		if err != nil {
			fmt.Println(tag, "error write conn: ", err)
			return err
		}
	} else if len(is.Options) >= 2 && is.Options[0] == "post" && is.Options[1] == "all" {
		var newClient client = client{}
		err := json.Unmarshal([]byte(is.Extras), &newClient)
		if err != nil {
			fmt.Println(tag, "error unmarshel for post all: ", err)
			return nil
		}
		cd.Clients[fmt.Sprint(t.Id)] = newClient

		writeFile()
		nets := server.NetStruct{
			Command: "Message",
			Options: []string{"receive post all command"},
		}
		err = t.Write(nets)
		if err != nil {
			fmt.Println(tag, "error write when receive post all")
			return err
		}

		for deviceIdStr, _ := range newClient.Devices {
			var deviceIdInt int
			deviceIdInt, err = strconv.Atoi(deviceIdStr)
			getDevice(uint(deviceIdInt))
		}
	}
	return nil
}

func writeFile() {
	file, err := os.OpenFile(clockFile, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(tag, "error open file for writeFile(): ", err)
		return
	}
	defer file.Close()
	b, err := json.Marshal(cd)
	if err != nil {
		fmt.Println(tag, "Error marshal for writeFile: ", err)
		return
	}
	_, err = file.Write(b)
	if err != nil {
		fmt.Println(tag, "Error WriteFile: ", err)
	}
	fmt.Println(tag, "successfully write to clock.json")
}

//get: 得到该闹钟的所有信息
func DealClockForDevice(is server.NetStruct, t server.TcpServer) error {
	if len(is.Options) >= 1 && is.Options[0] == "get" {
		getDevice(t.Id)
	}
	return nil
}

func getDeviceKeys(deviceId uint) (clientKey string, deviceKey string) {
	var client string
	var device string
LOOP:
	for k0, v := range cd.Clients {
		for k, _ := range v.Devices {
			if id, _ := strconv.ParseUint(k, 10, 32); uint(id) == deviceId {
				client = k0
				device = k
				break LOOP
			}
		}
	}
	return client, device
}

//直接发送给目标机器
func getDevice(deviceId uint) {
	clientKey, deviceKey := getDeviceKeys(deviceId)
	var content device = cd.Clients[clientKey].Devices[deviceKey]
	writeByte, err := json.Marshal(content)
	if err != nil {
		fmt.Println(tag, "error marshel: ", err)
	}

	nets := server.NetStruct{
		Command: "device",
		Options: []string{"total"},
		Extras:  string(writeByte),
	}

	server.WriteToDevice(deviceId, nets)
}
