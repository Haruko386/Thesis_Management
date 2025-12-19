package views

import (
	"ThesisManagement/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func LoginHandler(c *gin.Context) {
	// 初始化请求
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// 查询用户数据是否存在
	var user models.User
	err := models.DB.Where("username = ? ", req.Username).
		First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "账号或密码错误"})
		return
	}
	// 验证密码
	if !checkPassword(user.Password, req.Password) {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "账号或密码错误"})
		return
	}
	println(1)
	// 登录成功：设置一个名为 "user_session" 的 Cookie，有效期 1 小时
	c.SetCookie("user_session", user.Username, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
		"user": gin.H{
			"username": user.Username,
		},
	})
}
