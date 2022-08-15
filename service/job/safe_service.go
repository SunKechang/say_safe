package job

import (
	"fmt"
	"gin-test/database/safe"
	"gin-test/util/lock"
	"gin-test/util/log"
	"sync"
)

//一个service对应一种任务的执行
type SafeJobService struct {
	wg        sync.WaitGroup
	reentLock lock.ReentrantLock
}

func NewSafeJobService() *SafeJobService {
	return &SafeJobService{}
}

func (p *SafeJobService) Stop() {
	// TODO wait all jobs done
}

func (p *SafeJobService) ExecJobs() {
	if !p.reentLock.Lock() {
		log.Log(fmt.Sprintf("[job], reentrant locked\n"))
		return
	}

	defer p.reentLock.Free()

	safeJobDao := safe.NewSafeJobDao()
	jobs, err := safeJobDao.GetExistingJob()
	if err != nil {
		log.Log(fmt.Sprintf("safe job execjobs failed: %s\n", err))
	}
	for _, j := range jobs {
		p.execOneJob(j)
	}
	p.wg.Wait()
}

func (p *SafeJobService) execOneJob(jobInfo safe.SafeJobInfo) {
	var safeJob Job
	safeJob = NewCommitSafeJob(jobInfo)
	p.wg.Add(1)
	go func(job Job) {
		defer p.wg.Done()
		job.Exec()
	}(safeJob)
}
