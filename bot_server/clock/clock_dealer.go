//
package clock

import (
	"encoding/json"
	"first_server/server"
	"fmt"
	"io/ioutil"
	"os"
)

const tag = "clock "

type clock struct {
	Cron     string `json:"cron"`
	Text     string `json:"text"`
	IsActive bool   `json:"is_active"`
}

type device struct {
	//FixMe: 这里不知道可不可以用 uint
	ClientId uint    `json:"client_id"`
	Clocks   []clock `json:"clocks"`
}

type clockData struct {
	Devices map[string]device `json:"devices"`
}

var cd *clockData

const clockFile = "data/clock.json"

func init() {
	cd = &clockData{}

	file, err := os.Open(clockFile)
	if err != nil {
		fmt.Println(tag, "error open file: ", err)
	}
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
func DealClockForClient(is server.InputStruct, t server.TcpServer) error {
	if len(is.Options) >= 1 && is.Options[0] == "get" {
		//TODO:这里需要增加一层逻辑，只给当前 client 管理的 devices 的信息
		writeByte, err := json.Marshal(cd.Devices[fmt.Sprint(t.Id)])
		if err != nil {
			fmt.Println(tag, "error marshel: ", err)
		}
		writeByte = append(writeByte, byte('\n'))
		err = t.Write(writeByte)
		if err != nil {
			fmt.Println(tag, "error write conn: ", err)
			return err
		}
	}
	return nil
}
