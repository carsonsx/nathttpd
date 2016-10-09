package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type mqconfig struct {
	MQUrl  string `yaml:"mqserver"`
	MQName string `yaml:"mqname"`
}

var MQConf = mqconfig{"amqp://localhost:5672/", "http_nat_queue"}

func LoadConf(path string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	yaml.Unmarshal(c, &MQConf)
}

func init() {
	LoadConf("mq.yml")
}
