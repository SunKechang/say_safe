package safe

import (
	"gin-test/database"
	"gin-test/util"
	"strings"
	"time"
)

const (
	SafeJobPrefix = "sj"
)

type SafeJob struct {
	PK         int64     `json:"pk" gorm:"column:pk;primary_key;AUTO_INCREMENT"`
	ID         string    `json:"id" gorm:"column:id;"`
	UserId     string    `json:"user_id" gorm:"column:user_id;"`
	Path       string    `json:"path" gorm:"column:path;"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;default:"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:"`
	IsDelete   int       `json:"is_delete" gorm:"column:is_delete;default:"`
}

type SafeJobInfo struct {
	ID       string `json:"id" gorm:"column:id;"`
	UserId   string `json:"user_id" gorm:"column:user_id;"`
	Password string `json:"password" gorm:"column:password;"`
	Path     string `json:"path" gorm:"column:path;"`
}

func (u SafeJob) TableName() string {
	return database.SafeJobTableName
}

type SafeJobDao struct {
	database.BaseDao
}

func NewSafeJobDao() *SafeJobDao {
	return &SafeJobDao{
		database.BaseDao{Engine: database.GetDB()},
	}
}

func (p *SafeJobDao) CreateSafeJob(job *SafeJob) error {
	job.ID = p.GetUUID()
	q := p.GetDB()
	q.Create(job)
	return q.Error
}

func (p *SafeJobDao) GetJobByUserID(userId string) (*SafeJob, int64, error, error) {
	q := p.GetDB().Table(database.SafeJobTableName)
	res := &SafeJob{}
	q = q.Where("user_id = ?", userId)
	q.Where("is_delete = ?", database.Undeleted).First(&res)

	count := int64(0)
	q2 := p.GetDB().Table(database.SafeLogTableName)
	q2 = q2.Where("user_id = ?", userId)
	q2 = q2.Where("is_delete = ?", database.Undeleted)
	q2.Count(&count)
	return res, count, q.Error, q2.Error
}

func (p *SafeJobDao) GetExistingJob() ([]SafeJobInfo, error) {
	q := p.GetDB()
	q = q.Model(&SafeJob{})

	selects := []string{"safe_job.id as id",
		"safe_job.user_id as user_id",
		"user.password as password",
		"safe_job.path as path"}
	q = q.Select(strings.Join(selects, ","))
	q = q.Joins("left join user on user.id = safe_job.user_id")
	q = q.Where("safe_job.is_delete = ?", database.Undeleted)
	res := make([]SafeJobInfo, 0)
	q.Scan(&res)
	return res, q.Error
}

func (p *SafeJobDao) DeleteJobsByUserID(userId string) error {
	q := p.GetDB()
	q = q.Model(&SafeJob{})
	q = q.Where("user_id = ?", userId)
	q.Update("is_delete", database.Deleted)
	return q.Error
}

// GetUUID 获取uuid
func (p *SafeJobDao) GetUUID() string {
	return util.NewShortIDString(SafeJobPrefix)
}
