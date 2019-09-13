package common

import "encoding/json"

type Job struct {
	Name     string `json:""` //Task name
	Command  string `json:""` //Shell Command
	CronExpr string `json:""` //Cron Expression
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errno int, msg string, data interface{}) (response []byte, err error) {
	var (
		rsp Response
	)

	rsp.Errno = errno
	rsp.Msg = msg
	rsp.Data = data

	return json.Marshal(rsp)
}
