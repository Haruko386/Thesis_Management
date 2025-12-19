package views

import (
	"ThesisManagement/models"
	"github.com/gen2brain/go-fitz"
	"github.com/gin-gonic/gin"
	"image/jpeg"
	"net/http"
	"path/filepath"
	"time"
)

type thesisReq struct {
	ThesisName string `json:"thesis_name" binding:"required"`
}

func GetThesis(c *gin.Context) {
	// 获取请求
	var req thesisReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//
}

func GetPaperCover(c *gin.Context) {
	fileName := "ApDepth.pdf"
	filePath := filepath.Join("assets", "Paper", fileName)

	// 打开文件
	doc, err := fitz.New(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取PDF文件"})
		return
	}
	defer func(doc *fitz.Document) {
		err := doc.Close()
		if err != nil {
			panic(err)
		}
	}(doc)

	//提取第一页作为图片
	img, err := doc.Image(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转换图片失败"})
		return
	}

	// 设置响应头为图片格式并返回流
	c.Header("Content-Type", "image/jpeg")
	// 将图像以 JPEG 格式写入响应，Quality 80 可以在清晰度和体积间取得平衡
	_ = jpeg.Encode(c.Writer, img, &jpeg.Options{Quality: 80})
}

func PostUploadHandler(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	classification := c.PostForm("classification")
	//publicDate := c.PostForm("public_date")

	file, _ := c.FormFile("file")
	dst := filepath.Join("assets", "paper", file.Filename)
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		panic(err)
	}

	newThesis := models.Thesis{
		Title:          title,
		Author:         author,
		Classification: classification,
		StoragePath:    file.Filename,
		PublicDate:     time.Now(),
	}
	models.DB.Create(&newThesis)
	c.Redirect(http.StatusMovedPermanently, "/")
}
