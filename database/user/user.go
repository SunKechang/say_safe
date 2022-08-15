package user

import (
	"gin-test/database"
	"gin-test/util"
	"time"
)

const (
	UserPrefix = "u"
)

type User struct {
	PK         int64     `json:"pk" gorm:"column:pk;primary_key;AUTO_INCREMENT"`
	ID         string    `json:"id" gorm:"column:id;"`
	UserName   string    `json:"user_name" gorm:"column:user_name;"`
	Salt       string    `json:"salt" gorm:"column:salt;"`
	Password   string    `json:"password" gorm:"column:password;"`
	Class      string    `json:"class" gorm:"column:klass;"`
	IsMan      int       `json:"is_man" gorm:"column:is_man;"`
	EndTime    time.Time `json:"end_time" gorm:"column:end_time;default:"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;default:"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:"`
	IsDelete   int       `json:"is_delete" gorm:"column:is_delete;default:"`
}

func (u User) TableName() string {
	return database.UserTableName
}

type UserDao struct {
	database.BaseDao
}

func NewUserDao() *UserDao {
	return &UserDao{
		database.BaseDao{Engine: database.GetDB()},
	}
}

func (p *UserDao) CreateUser(u *User) error {
	q := p.GetDB()
	if u.ID == "" {
		u.ID = p.GetUUID()
	}

	q = q.Create(u)
	return q.Error
}

func (p *UserDao) GetUserByID(id string) (*User, error) {
	q := p.GetDB().Table(database.UserTableName)
	res := &User{}
	q = q.Where("id = ?", id)
	q.Where("is_delete = ?", database.Undeleted).First(&res)
	return res, q.Error
}

// GetUUID 获取uuid
func (p *UserDao) GetUUID() string {
	return util.NewShortIDString(UserPrefix)
}
