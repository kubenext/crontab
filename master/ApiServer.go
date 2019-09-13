package master

import (
	"encoding/json"
	"fmt"
	"github.com/kubenext/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	G_apiServer *ApiServer
)

// Http interface
type ApiServer struct {
	httpServer *http.Server
}

// Save job interface
func handleJobSave(response http.ResponseWriter, request *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	if err = request.ParseForm(); err != nil {
		goto ERR
	}

	postJob = request.PostForm.Get("job")

	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	if oldJob, err = G_jobMgr.Save(&job); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		response.Write(bytes)
	}

	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		response.Write(bytes)
	}
}

func handleJobDelete(response http.ResponseWriter, request *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	if err = request.ParseForm(); err != nil {
		goto ERR
	}

	name = request.PostForm.Get("name")

	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		response.Write(bytes)
	}

	return

ERR:

	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		response.Write(bytes)
	}

}

func handleJobList(response http.ResponseWriter, request *http.Request) {
	var (
		jobs  []*common.Job
		err   error
		bytes []byte
	)
	if jobs, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", jobs); err == nil {
		response.Write(bytes)
	}

	return

ERR:

	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		response.Write(bytes)
	}

}

func handleJobKill(response http.ResponseWriter, request *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	if err = request.ParseForm(); err != nil {
		goto ERR
	}

	name = request.PostForm.Get("name")

	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		response.Write(bytes)
	}

	return

ERR:

	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		response.Write(bytes)
	}

}

// Initialization service
func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	staticDir = http.Dir(G_config.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler))

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ServerPort)); err != nil {
		return err
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ServerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ServerWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	go func() {
		if err := httpServer.Serve(listener); err != nil {
			fmt.Println(err)
		}
	}()

	return
}
