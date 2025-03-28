package service

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FullTextSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.ErrorResponse(c, "Vui lòng cung cấp từ khóa tìm kiếm", http.StatusBadRequest, nil)
		return
	}

	db := config.GetDB()
	var files []models.File
	likeQuery := "%" + keyword + "%"
	query := fmt.Sprintf("MATCH(file_name, content) AGAINST('%s' IN BOOLEAN MODE)", keyword)
	err := db.Where(query).Or("file_name LIKE ? OR content LIKE ?", likeQuery, likeQuery).Find(&files).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorResponse(c, "Lỗi tìm kiếm", http.StatusInternalServerError, nil)
		return
	}
	utils.SuccessResponse(c, "Tìm kiếm thành công", http.StatusOK, files)

}
