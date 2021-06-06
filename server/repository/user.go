package repository

import (
	"github.com/wxmsummer/AirConditioner/server/model"
	"gorm.io/gorm"
)

// 用户数据库操作相关接口
type UserRepository interface {
	FindById(int) (*model.User, error)                        // 根据Id查找用户
	Create(*model.User) error                                 // 创建用户
	Update(*model.User) (*model.User, error)                  // 更新用户
	FindByField(string, string, string) ([]model.User, error) // 条件查询
	FindAllUsers() (users []model.User, err error)            // 查询所有用户
}

type UserOrm struct {
	Db *gorm.DB
}

func (orm *UserOrm) FindById(id int) (user model.User, err error) {
	err = orm.Db.Select("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (orm *UserOrm) FindByPhone(phone string) (user model.User, err error) {
	err = orm.Db.Select("phone = ?", phone).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (orm *UserOrm) Create(user *model.User) error {
	err := orm.Db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *UserOrm) Update(user *model.User) error {
	err := orm.Db.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (orm *UserOrm) FindByField(key, value, fields string) (users []model.User, err error) {
	if len(fields) == 0 {
		fields = "*"
	}
	err = orm.Db.Select(fields).Where(key+" = ?", value).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 查询所有用户
func (orm *UserOrm) FindAllUsers() (users []model.User, err error) {
	err = orm.Db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
