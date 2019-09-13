package master

import (
	"context"
	"encoding/json"
	"github.com/kubenext/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	G_jobMgr *JobMgr
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return err
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

func (mgr *JobMgr) Save(job *common.Job) (oldJob *common.Job, err error) {
	var (
		jobKey      string
		jobValue    []byte
		putResponse *clientv3.PutResponse
		oldJobObj   common.Job
	)

	jobKey = common.JOB_SAVE_DIR + job.Name

	if jobValue, err = json.Marshal(job); err != nil {
		return nil, err
	}

	if putResponse, err = mgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return nil, err
	}

	if putResponse.PrevKv != nil {
		if err = json.Unmarshal(putResponse.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}

		oldJob = &oldJobObj
	}

	return

}

func (mgr *JobMgr) DeleteJob(name string) (oldJob *common.Job, err error) {
	var (
		jobKey      string
		delResponse *clientv3.DeleteResponse
		oldJobObj   common.Job
	)
	jobKey = common.JOB_SAVE_DIR + name

	if delResponse, err = mgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delResponse.PrevKvs) != 0 {
		if err = json.Unmarshal(delResponse.PrevKvs[0].Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return
}
