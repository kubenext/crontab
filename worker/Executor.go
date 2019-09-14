package worker

import (
	"context"
	"github.com/kubenext/crontab/common"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)

		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
		}

		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		result.StartTime = time.Now()

		err = jobLock.TryLock()
		defer jobLock.Unlock()

		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			result.StartTime = time.Now()
			cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)

			output, err = cmd.CombinedOutput()

			result.EndTime = time.Now()
			result.Output = output
			result.Err = err

		}

		G_scheduler.PushJobResult(result)

	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}