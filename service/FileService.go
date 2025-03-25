package service

import (
	"backend/config"
	"backend/utils"
	"fmt"
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
	file, err := c.FormFile("file") // lấy file từ request
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

	// Lấy phần mở rộng của file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		utils.ErrorResponse(c, "Định dạng file không được phép. Chỉ hỗ trợ pdf, png, jpg, jpeg, txt", http.StatusBadRequest, nil)
		return
	}

	// Lấy đường dẫn thư mục từ đường dẫn, nếu không có thì mặc định là uploads
	userDir := c.PostForm("path")
	if userDir == "" {
		userDir = config.AppConfig.Storage.UploadPath
	}

	userDir = filepath.Clean(userDir)
	if strings.Contains(userDir, "..") { // ngăn chặn đường dẫn không hợp lệ
		utils.ErrorResponse(c, "Thư mục không hợp lệ", http.StatusBadRequest, nil)
		return
	}
	if ok, err := utils.CheckFolder(config.AppConfig.Storage.UploadPath); !ok {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể tạo thư mục: %s", err.Error()), http.StatusInternalServerError, nil)
	}
	//Xây dựng đường dẫn lưu file
	filePath := filepath.Join(userDir, file.Filename)

	//Lưu file vào đường dẫn chỉ định
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể lưu file: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	// Thông báo trả về kết quả nếu upload thành công
	utils.SuccessResponse(c, "Tải lên thành công", http.StatusOK, gin.H{"path": filePath})
}
