package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	G_config *Config
)

type Config struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

// Load configuration
func InitConfig(fileName string) (err error) {
	var (
		content []byte
		config  Config
	)

	if content, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if err = json.Unmarshal(content, &config); err != nil {
		return
	}

	G_config = &config

	fmt.Println(config)

	return
}
