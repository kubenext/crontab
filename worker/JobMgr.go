package worker

import (
	"context"
	"github.com/kubenext/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

var (
	G_jobMgr *JobMgr
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
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
	watcher = clientv3.NewWatcher(client)

	G_jobMgr = &JobMgr{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	G_jobMgr.watchJobs()

	// 监听Killer
	G_jobMgr.watchKiller()

	return
}

func (mgr *JobMgr) watchKiller() (err error) {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
		job        *common.Job
	)

	watchChan = mgr.watcher.Watch(context.TODO(), common.JOB_KILL_DIR, clientv3.WithPrefix())

	for watchResp = range watchChan {
		for _, watchEvent = range watchResp.Events {
			switch watchEvent.Type {
			case mvccpb.PUT:
				jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
				job = &common.Job{Name: jobName}
				jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)
				G_scheduler.PushJobEvent(jobEvent)
			case mvccpb.DELETE:

			}
		}
	}

	return
}

func (mgr *JobMgr) watchJobs() (err error) {

	var (
		getResponse        *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResponse      clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	if getResponse, err = mgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvpair = range getResponse.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	go func() {
		watchStartRevision = getResponse.Header.Revision + 1
		watchChan = mgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())

		for watchResponse = range watchChan {
			for _, watchEvent = range watchResponse.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
				}

				G_scheduler.PushJobEvent(jobEvent)
			}
		}
	}()

	return
}

func (mgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, mgr.kv, mgr.lease)
	return
}
