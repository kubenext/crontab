package master

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	G_config *Config
)

type Config struct {
	ServerPort            int      `json:"serverPort"`
	ServerReadTimeout     int      `json:"serverReadTimeout"`
	ServerWriteTimeout    int      `json:"serverWriteTimeout"`
	EtcdEndpoints         []string `json:"etcdEndpoints"`
	EtcdDialTimeout       int      `json:"etcdDialTimeout"`
	Webroot               string   `json:"webroot"`
	MongodbUri            string   `json:"mongodbUri"`
	MongodbConnectTimeout int      `json:"mongodbConnectTimeout"`
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
