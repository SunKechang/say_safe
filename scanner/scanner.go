package scanner

import (
	"gin-test/service/job"
	"github.com/robfig/cron"
)

//代表各种任务的指挥者，可以按cron定时执行这一类型的任务
type Scanner struct {
	cron        *cron.Cron
	jobServices []job.JobService
}

func NewScanner() *Scanner {
	return &Scanner{
		cron:        cron.New(),
		jobServices: make([]job.JobService, 0),
	}
}

func (p *Scanner) AddService(interval string, service job.JobService) {
	p.jobServices = append(p.jobServices, service)
	p.cron.AddFunc(interval, func() {
		service.ExecJobs()
	})
}

func (s *Scanner) Start() (err error) {
	s.cron.Start()
	return
}

func (s *Scanner) Stop() (err error) {
	s.cron.Stop()
	for _, s := range s.jobServices {
		s.Stop()
	}
	return
}
