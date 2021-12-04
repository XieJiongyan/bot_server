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
	Cron string `json:"cron"`
	Text string `json:"text"`
}

type clockData struct {
	Clock map[string][]clock `json:"clock"`
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

func DealClock(is server.InputStruct, t server.TcpServer) error {
	if len(is.Options) >= 1 && is.Options[1] == "get" {
		writeByte, err := json.Marshal(cd.Clock[fmt.Sprint(t.Id)])
		if err != nil {
			fmt.Println(tag, "error marshel: ", err)
		}

		err = t.Write(writeByte)
		if err != nil {
			fmt.Println(tag, "error write conn: ", err)
			return err
		}
	}
	return nil
}
