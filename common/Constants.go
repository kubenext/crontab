package common

const (
	JOB_SAVE_DIR = "/cron/jobs/"
	JOB_KILL_DIR = "/cron/killer/"

	JOB_EVENT_SAVE   = 1
	JOB_EVENT_DELETE = 2
	JOB_EVENT_KILL   = 3

	JOB_LOCK_KEY = "/cron/lock/"

	JOB_WORKER_DIR = "/cron/workers/"
)
