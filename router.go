package main

import (
	"ThesisManagement/models"
	"ThesisManagement/views"
	"github.com/gin-gonic/gin"
	"net/http"
)

/* 业务逻辑 */
func initRoutes(r *gin.Engine) {
	// 127.0.0.1:8080
	/* 页面渲染 */
	r.GET("/", func(c *gin.Context) {
		// 判断是否登录(有无session返回)
		username, err := c.Cookie("user_session")
		println(c.Cookie("user_session"))
		isLogin := err == nil

		// 获取全部论文
		var theses []models.Thesis
		models.DB.Find(&theses)

		var categories []string
		models.DB.Model(&models.Thesis{}).Pluck("DISTINCT classification", &categories)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"isLogin":    isLogin,
			"user":       username,
			"theses":     theses,
			"categories": categories,
		})
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", nil)
	})
	r.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.tmpl", nil)
	})

	r.GET("/paper/cover", views.GetPaperCover)

	r.GET("/thesis", func(c *gin.Context) {
		c.HTML(http.StatusOK, "thesis.tmpl", nil)
	})

	r.GET("/logout", func(c *gin.Context) {
		c.SetCookie("user_session", "", -1, "/", "localhost", false, true)
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	/* api组件 */
	api := r.Group("/api")
	{
		api.POST("login", views.LoginHandler)
		api.GET("search", views.SearchHandler) // 占位符
		api.POST("/upload", views.PostUploadHandler)
	}
}
