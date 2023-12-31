package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/xuhaidong1/webook/internal/domain"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrNotFound           = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	//存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	//保持链路 with context
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflicts uint16 = 1062
		if mysqlErr.Number == uniqueConflicts {
			//邮箱冲突，因为user表只有一个唯一索引
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).First(&u, id).Error
	return u, err
}

func (dao *UserDAO) UpdateProfileById(ctx context.Context, u domain.User) error {
	user := &User{
		Id:           u.Id,
		Nickname:     u.Nickname,
		Birth:        u.Birth,
		Introduction: u.Introduction,
	}
	err := dao.db.WithContext(ctx).Model(&user).Updates(map[string]any{
		"nickname":     user.Nickname,
		"birth":        user.Birth,
		"introduction": user.Introduction,
	}).Error
	return err
}

// User 直接对应表结构
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	//全部用户唯一
	Email        string `gorm:"unique"`
	Password     string
	Nickname     string `gorm:"type:varchar(50)"`
	Birth        string
	Introduction string `gorm:"type:varchar(1000)"`
	//创建时间 毫秒数 使用UTC，避免时区问题
	Ctime int64
	//更新时间 毫秒数
	Utime int64
}
