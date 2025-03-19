package service

import (
	"backend/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		fmt.Println("Lỗi:", err)
		return
	}
	var allFile []string
	for _, file := range files {
		allFile = append(allFile, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{
		"file":        allFile,
		"status_code": http.StatusOK,
	})
}

func DownloadFile(c *gin.Context) {
	fileName := c.Param("filename") // Lấy tên file từ request
	filePath := filepath.Join(config.AppConfig.Storage.UploadPath, fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) { // Kiểm tra file có tồn tại hay không
		c.JSON(http.StatusNotFound, gin.H{"error": "File không tồn tại"})
		return
	}

	// Cấu hình header để trình duyệt hiểu rằng đây là file cần tải về
	c.Header("Content-Disposition", "attachment; filename="+fileName) //
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
func SearchFile(c *gin.Context) {
	keyword := c.Query("keyword") // Lấy từ khóa tìm kiếm từ query string

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vui lòng cung cấp từ khóa tìm kiếm"})
		return
	}

	files, err := os.ReadDir("uploads") // Đọc danh sách file trong thư mục "uploads"
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi khi đọc thư mục"})
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
		c.JSON(http.StatusNotFound, gin.H{
			"message":     "Không tìm thấy file nào",
			"status_code": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Danh sách file trùng với từ khóa",
		"files":       matchedFiles,
		"status_code": http.StatusOK,
	})
}

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file") // lấy file từ request
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ // Thông báo nếu không thể lấy file
			"error":       "Không thể lấy file",
			"status_code": http.StatusBadRequest,
		})
		return
	}

	// Lấy đường dẫn thư mục từ đường dẫn, nếu không có thì mặc định là uploads
	userDir := c.PostForm("path")
	if userDir == "" {
		userDir = config.AppConfig.Storage.UploadPath
	}

	userDir = filepath.Clean(userDir)
	if strings.Contains(userDir, "..") { // ngăn chặn đường dẫn không hợp lệ
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Thư mục không hợp lệ",
			"status_code": http.StatusBadRequest,
		})
		return
	}
	if _, err := os.Stat(userDir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
				log.Println("Lỗi tạo thư mục:", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       fmt.Sprintf("Không thể tạo thư mục: %v", err),
					"status_code": http.StatusInternalServerError,
				})
				return
			}
		} else {
			log.Println("Lỗi kiểm tra thư mục:", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":       fmt.Sprintf("Lỗi kiểm tra thư mục: %v", err),
				"status_code": http.StatusInternalServerError,
			})
			return
		}
	}

	//Xây dựng đường dẫn lưu file
	filePath := filepath.Join(userDir, file.Filename)

	//Lưu file vào đường dẫn chỉ định
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ // Trả về mã lỗi nếu không lưu được file
			"error":       "Không thể lưu file",
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	// Thông báo trả về kết quả nếu upload thành công
	c.JSON(http.StatusOK, gin.H{
		"message":     "Tải lên thành công",
		"path":        filePath,
		"status_code": http.StatusOK,
	})
}
