package main

import (
	"flag"
	"fmt"
	"github.com/kubenext/crontab/worker"
	"runtime"
	"time"
)

var (
	configFilePath string
)

func initArgs() {
	flag.StringVar(&configFilePath, "config", "./worker.json", "config path")
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

	if err = worker.InitConfig(configFilePath); err != nil {
		goto ERR
	}

	if err = worker.InitRegister(); err != nil {
		goto ERR
	}

	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

	return

ERR:
	fmt.Println(err)

}
