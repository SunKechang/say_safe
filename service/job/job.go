package job

type JobService interface {
	Stop()
	ExecJobs()
}

type Job interface {
	Exec()
}
