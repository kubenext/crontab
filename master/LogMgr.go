package master

import (
	"context"
	"github.com/kubenext/crontab/common"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	G_logMgr *LogMgr
)

func InitLogMgr() (err error) {
	var (
		client     *mongo.Client
		clientOpts *options.ClientOptions
		timeout    time.Duration
	)
	timeout = time.Duration(G_config.MongodbConnectTimeout) * time.Millisecond
	clientOpts = options.Client().ApplyURI(G_config.MongodbUri)
	clientOpts.ConnectTimeout = &timeout
	if client, err = mongo.Connect(context.TODO(), clientOpts); err != nil {
		return
	}

	G_logMgr = &LogMgr{
		client:        client,
		logCollection: client.Database("cron").Collection("log"),
	}

	return
}

func (logMgr *LogMgr) ListLog(name string, start int64, limit int64) (logArr []*common.JobLog, err error) {
	var (
		filter *common.JobLogFilter
		sort   *common.SortLogByStartTIme
		opt    *options.FindOptions
		cursor *mongo.Cursor
		jobLog *common.JobLog
	)

	logArr = make([]*common.JobLog, 0)

	filter = &common.JobLogFilter{JobName: name}
	sort = &common.SortLogByStartTIme{SortOrder: -1}

	opt = options.Find()
	opt.SetSort(sort)
	opt.SetSkip(start)
	opt.SetLimit(limit)

	if cursor, err = logMgr.logCollection.Find(context.TODO(), filter, opt); err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}

		logArr = append(logArr, jobLog)
	}

	return
}
