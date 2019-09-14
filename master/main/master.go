package main

import (
	"flag"
	"fmt"
	"github.com/kubenext/crontab/master"
	"runtime"
	"time"
)

var (
	configFilePath string
)

func initArgs() {
	flag.StringVar(&configFilePath, "config", "./master.json", "config path")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	var (
		err error
	)

	initArgs()

	initEnv()

	if err = master.InitConfig(configFilePath); err != nil {
		goto ERR
	}

	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return

ERR:
	fmt.Println(err)

}
