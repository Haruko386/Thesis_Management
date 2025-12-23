package views

import (
	"ThesisManagement/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func EditThesis(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找
	var thesis models.Thesis
	if err := models.DB.First(&thesis, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取表单中的数据
	title := c.PostForm("title")
	author := c.PostForm("author")
	journal := c.PostForm("journal")
	classification := c.PostForm("classification")
	publicDateStr := c.PostForm("publicdate")
	publicDate, err := time.Parse("2006-01-02", publicDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "发布日期格式错误"})
		return
	}

	thesis.Title = title
	thesis.Author = author
	thesis.Journal = journal
	thesis.Classification = classification
	thesis.PublicDate = publicDate

	if err := models.DB.Save(&thesis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "论文更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "论文更新成功"})
}
