package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"github.com/carsonsx/nathttpd/const"
)

type MQConfig struct {
	Url      string `yaml:"url"`
	ReqQueue string `yaml:"request_queue"`
	ResQueue string `yaml:"response_queue"`
}

var MQConf = MQConfig{
	constant.DEFAULT_CONNECTION_URL,
	constant.DEFAULT_REQUEST_QUEUE,
	constant.DEFAULT_RESPONSE_QUEUE,
}

func LoadConf(path string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	yaml.Unmarshal(c, &MQConf)
}

func init() {
	LoadConf("mq.yml")
}
