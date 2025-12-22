package main

import (
	"ThesisManagement/models"
	"ThesisManagement/views"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type CategoryCount struct {
	Name  string
	Count int
}

func RequireLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("user_session")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequireLoginAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("user_session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "请先登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

/* 业务逻辑 */
func initRoutes(r *gin.Engine) {
	// 127.0.0.1:8080

	/* 页面渲染 */
	r.GET("/", func(c *gin.Context) {
		// 登录态
		username, err := c.Cookie("user_session")
		isLogin := err == nil

		// 分类过滤（可选：点分类时只看该分类）
		currentCls := c.Query("c") // 例如 /?c=计算机视觉
		var theses []models.Thesis
		if currentCls != "" && currentCls != "全部" {
			models.DB.
				Where("classification LIKE ?", "%"+currentCls+"%").
				Find(&theses)
		} else {
			models.DB.Order("created_at desc").Find(&theses)
			currentCls = "全部"
		}

		categories := make(map[string]int)

		var all []models.Thesis
		models.DB.Find(&all)

		for i := range all {
			if all[i].Classification == "" {
				continue
			}
			list := strings.Split(all[i].Classification, ";")
			for _, item := range list {
				categories[item]++
			}
		}

		total := len(all)

		// 给 upload.tmpl 用的“可选分类列表”（去掉“未分类”重复项）
		var categoryOptions []string
		models.DB.Model(&models.Thesis{}).Where("classification <> '' AND classification <> '未分类'").
			Pluck("DISTINCT classification", &categoryOptions)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"isLogin":         isLogin,
			"user":            username,
			"theses":          theses,
			"categories":      categories,
			"totalCount":      total,
			"currentCls":      currentCls,
			"categoryOptions": categoryOptions,
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", nil)
	})

	// 只有登录后才能打开上传页
	r.GET("/upload", RequireLoginPage(), func(c *gin.Context) {
		var categoryOptions []string
		models.DB.Model(&models.Thesis{}).Where("classification <> '' AND classification <> '未分类'").
			Pluck("DISTINCT classification", &categoryOptions)

		c.HTML(http.StatusOK, "upload.tmpl", gin.H{
			"categoryOptions": categoryOptions,
		})
	})

	// 论文封面：按论文ID生成第一页jpg
	r.GET("/paper/cover/:id", views.GetPaperCover)

	// 论文阅读页：按论文ID读取 StoragePath
	r.GET("/thesis/:id", func(c *gin.Context) {
		var thesis models.Thesis
		if err := models.DB.First(&thesis, c.Param("id")).Error; err != nil {
			c.String(http.StatusNotFound, "论文不存在")
			return
		}
		c.HTML(http.StatusOK, "thesis.tmpl", gin.H{
			"thesis": thesis,
		})
	})

	r.DELETE("/delete/:id", func(c *gin.Context) {
		var thesis models.Thesis
		if err := models.DB.First(&thesis, c.Param("id")).Error; err != nil {
			c.String(http.StatusNotFound, "删除失败")
			return
		}
	})

	// 退出登录
	r.GET("/logout", func(c *gin.Context) {
		// MaxAge=-1 删除 cookie；domain 传 "" 不写死
		c.SetCookie("user_session", "", -1, "/", "", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	/* api组件 */
	api := r.Group("/api")
	{
		api.POST("login", views.LoginHandler)
		api.GET("search", views.SearchHandler) // 占位符

		// 只有登录后才能上传
		authed := api.Group("")
		authed.Use(RequireLoginAPI())
		authed.POST("/upload", views.PostUploadHandler)
	}
}
