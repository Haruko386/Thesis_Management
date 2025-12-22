package main

import (
	"ThesisManagement/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func main() {
	// 初始化gin框架
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.Static("/asset", "./assets")

	// 初始化数据库
	err := models.InitMysql()
	if err != nil {
		return
	}
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(models.DB)
	//迁移
	models.DB.AutoMigrate(&models.User{}, &models.Thesis{})

	// 初始化路由
	initRoutes(r)

	// run
	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
