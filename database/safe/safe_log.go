package safe

import (
	"gin-test/database"
	"gin-test/util"
	"time"
)

const (
	SafeLogPrefix = "sl"
)

type SafeLog struct {
	PK         int64     `json:"pk" gorm:"column:pk;primary_key;AUTO_INCREMENT"`
	ID         string    `json:"id" gorm:"column:id;"`
	UserId     string    `json:"user_id" gorm:"column:user_id;"`
	JobId      string    `json:"job_id" gorm:"column:job_id;"`
	Success    int       `json:"success" gorm:"column:success;"`
	Result     string    `json:"result" gorm:"column:result;"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;default:"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:"`
	IsDelete   int       `json:"is_delete" gorm:"column:is_delete;default:"`
}

func (s SafeLog) TableName() string {
	return database.SafeLogTableName
}

type SafeLogDao struct {
	database.BaseDao
}

func NewSafeLogDao() *SafeLogDao {
	return &SafeLogDao{
		database.BaseDao{Engine: database.GetDB()},
	}
}

func (p *SafeLogDao) CreateLog(log *SafeLog) error {
	log.ID = p.GetUUID()
	q := p.GetDB()
	q.Create(log)
	return q.Error
}

// GetUUID 获取uuid
func (p *SafeLogDao) GetUUID() string {
	return util.NewShortIDString(SafeLogPrefix)
}

func (p *SafeLogDao) GetSafeLog(userId string, pageNo, pageSize int) ([]SafeLog, error) {
	q := p.GetDB().Table(database.SafeLogTableName)
	q = q.Where("user_id = ?", userId)
	q = q.Where("is_delete = ?", database.Undeleted)
	offset := (pageNo - 1) * pageSize
	q = q.Offset(offset).Limit(pageSize)
	res := make([]SafeLog, 0)
	q.Scan(&res)
	return res, q.Error
}
