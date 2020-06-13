package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取页码
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * GetPageSize(c)
	}
	return result
}

// 获取每页的size
func GetPageSize(c *gin.Context) int {
	PageSize := 20
	limit, _ := com.StrTo(c.Query("limit")).Int()
	if limit > 1 {
		PageSize = limit
	}
	return PageSize
}
