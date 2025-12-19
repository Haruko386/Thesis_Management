package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func initMysql() (err error) {
	dsn := "root:password@(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.DB().Ping()
	return err
}

// 数据库

// User /* 只有我自己一个人，我只要一个用户名+密码验证登录就行了，其它目前来看没啥必要 */
type User struct {
	gorm.Model
	Username string
	Password string
}

// Thesis /* 论文数据；这部分只是给前端读取的，用于展示论文而已，看论文我们会直接读取pdf，因此这部分我们只需要少部分信息即可 */
type Thesis struct {
	gorm.Model
	Title       string
	Author      string
	StoragePath string
}
