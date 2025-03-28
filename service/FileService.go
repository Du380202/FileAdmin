package service

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	var allFile []models.File
	db := config.GetDB()
	db.Raw("SELECT * FROM files").Scan(&allFile)

	c.JSON(http.StatusOK, gin.H{
		"file":        allFile,
		"status_code": http.StatusOK,
	})
}

func DownloadFile(c *gin.Context) {
	fileName := c.Param("filename") // Lấy tên file từ request
	pathSearch := c.Param("path")
	filePath := filepath.Join(pathSearch, fileName)
	if _, err := os.Stat(pathSearch); os.IsNotExist(err) {
		utils.ErrorResponse(c, "Đường dẫn không tồn tại", http.StatusBadRequest, nil)
		return
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // Kiểm tra file có tồn tại hay không
		utils.ErrorResponse(c, "File không tồn tại", http.StatusNotFound, nil)
		return
	}

	// Cấu hình header để trình duyệt hiểu rằng đây là file cần tải về
	c.Header("Content-Disposition", "attachment; filename="+fileName) //
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}

func SearchFile(c *gin.Context) {
	keyword := c.Query("keyword") // Lấy từ khóa tìm kiếm từ query string
	pathSearch := c.Query("path")
	if keyword == "" {
		utils.ErrorResponse(c, "Vui lòng cung cấp từ khóa tìm kiếm", http.StatusBadRequest, nil)
		return
	}

	if pathSearch == "" {
		utils.ErrorResponse(c, "Vui lòng cung cấp đường dẫn tìm kiếm", http.StatusBadRequest, nil)
		return
	}

	if _, err := os.Stat(pathSearch); os.IsNotExist(err) {
		utils.ErrorResponse(c, "Đường dẫn tìm kiếm không tồn tại", http.StatusBadRequest, nil)
		return
	}

	files, err := os.ReadDir(pathSearch) // Đọc danh sách file trong thư mục do người dùng cung cấp
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi khi đọc thư mục: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	var matchedFiles []string
	// Lặp qua danh sách file để tìm file có chứa từ khóa
	for _, file := range files {
		if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(keyword)) {
			matchedFiles = append(matchedFiles, file.Name())
		}
	}

	if len(matchedFiles) == 0 {
		utils.ErrorResponse(c, "Không tìm thấy file nào", http.StatusNotFound, nil)
		return
	}

	utils.SuccessResponse(c, "Danh sách file trùng với từ khóa", http.StatusOK, gin.H{"files": matchedFiles})
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể lấy file: %s", err.Error()), http.StatusBadRequest, nil)
		return
	}

	allowedExtensions := map[string]bool{
		".pdf":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".txt":  true,
		".doc":  true,
		".docx": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		utils.ErrorResponse(c, "Định dạng file không được phép. Chỉ hỗ trợ pdf, png, jpg, jpeg, txt, doc, docx", http.StatusBadRequest, nil)
		return
	}

	userDir := c.PostForm("path")
	content := c.PostForm("content")
	if userDir == "" {
		userDir = config.AppConfig.Storage.UploadPath
	}

	userDir = filepath.Clean(userDir)
	if strings.Contains(userDir, "..") {
		utils.ErrorResponse(c, "Thư mục không hợp lệ", http.StatusBadRequest, nil)
		return
	}

	if ok, err := utils.CheckFolder(userDir); !ok {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể tạo thư mục: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	filePath := filepath.Join(userDir, file.Filename)
	if _, err := os.Stat(filePath); err == nil {
		timestamp := time.Now().Format("20060102_150405")
		newFileName := fmt.Sprintf("%s_%s%s", strings.TrimSuffix(file.Filename, ext), timestamp, ext)
		filePath = filepath.Join(userDir, newFileName)
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể lưu file: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	// userID, exists := c.Get("userID") // Giả sử userID được lấy từ middleware xác thực
	// if !exists {
	// 	utils.ErrorResponse(c, "Không tìm thấy thông tin người dùng", http.StatusUnauthorized, nil)
	// 	return
	// }

	fileRecord := models.File{
		FileName: file.Filename,
		FilePath: filePath,
		Content:  content, // Nếu cần lưu nội dung file, có thể đọc file và lưu vào đây
		UserID:   1,
	}

	db := config.GetDB()

	if err := db.Create(&fileRecord).Error; err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể lưu vào database: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	utils.SuccessResponse(c, "Tải lên thành công", http.StatusOK, gin.H{"file": fileRecord})
}
