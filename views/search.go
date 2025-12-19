package views

import "github.com/gin-gonic/gin"

type searchReq struct {
	Title      string `json:"title"`
	PublicDate string `json:"publicDate"`
}

func SearchHandler(c *gin.Context) {

}
