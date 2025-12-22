package views

import (
	"ThesisManagement/models"
	"github.com/gen2brain/go-fitz"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 获取论文封面（从PDF中提取第一页）
func GetPaperCover(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var thesis models.Thesis
	if err := models.DB.First(&thesis, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "论文不存在"})
		return
	}

	// StoragePath 形如 Paper/xxx.pdf
	filePath := filepath.Join("assets", thesis.StoragePath)

	doc, err := fitz.New(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取PDF文件"})
		return
	}
	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换图片失败"})
		return
	}

	c.Header("Content-Type", "image/jpeg")
	_ = jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 80})
}

func DeletePaper(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var thesis models.Thesis
	if err := models.DB.First(&thesis, id).Error; err != nil {
		c.JSON(http.StatusNotFound, "查找论文失败")
		return
	}
	if err := models.DB.Delete(&thesis).Error; err != nil {
		c.JSON(http.StatusForbidden, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// 处理上传论文
func PostUploadHandler(c *gin.Context) {
	title := strings.TrimSpace(c.PostForm("title"))
	author := strings.TrimSpace(c.PostForm("author"))
	classification := strings.TrimSpace(c.PostForm("classification"))
	journal := strings.TrimSpace(c.PostForm("journal"))
	if classification == "" {
		classification = "未分类"
	}

	// 处理多个分类，按 ; 分割
	classifications := strings.Split(classification, ";")

	// 处理发布日期
	publicDateStr := c.PostForm("publicdate")
	publicDate, err := time.Parse("2006-01-02", publicDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "发布日期格式错误"})
		return
	}

	// 文件处理
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "未选择文件"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "仅支持上传 PDF"})
		return
	}

	// 确保目录存在：assets/Paper
	saveDir := filepath.Join("assets", "Paper")
	_ = os.MkdirAll(saveDir, 0755)

	// 防止重名：时间戳_原文件名
	base := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	base = strings.ReplaceAll(base, " ", "_")
	newName := strconv.FormatInt(time.Now().Unix(), 10) + "_" + base + ext

	dst := filepath.Join(saveDir, newName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存文件失败"})
		return
	}

	// 将分类字段处理成单个字符串
	classificationStr := strings.Join(classifications, ";")

	// 创建新论文记录
	newThesis := models.Thesis{
		Title:          title,
		Author:         author,
		Journal:        journal,
		Classification: classificationStr,                                 // 将多个分类存为一个字符串
		StoragePath:    filepath.ToSlash(filepath.Join("Paper", newName)), // Paper/xxx.pdf
		PublicDate:     publicDate,                                        // 使用传递的 publicDate
	}

	// 保存论文记录到数据库
	if err := models.DB.Create(&newThesis).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "写入数据库失败"})
		return
	}

	// 上传成功，重定向到首页
	c.Redirect(http.StatusFound, "/")
}
