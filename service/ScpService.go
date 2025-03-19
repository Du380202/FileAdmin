package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"backend/config"

	"github.com/gin-gonic/gin"
)

// Ghi log lỗi SCP vào file error.log
func logError(message string, err error, output string) {
	f, _ := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	logger := log.New(f, "SCP_ERROR: ", log.LstdFlags)
	logger.Println(message, err.Error(), output)
}

// Hàm thực hiện lệnh SCP để chuyển file lên máy chủ từ xa
func scpFile(localFile, remoteFile string) error {
	cmd := exec.Command("scp", "-P", config.AppConfig.SCP.RemotePort, localFile, // Tạo lệnh SCP
		fmt.Sprintf("%s@%s:%s", config.AppConfig.SCP.Username, config.AppConfig.SCP.RemoteHost, remoteFile))

	// Thực thi lệnh SCP và lấy đầu ra (bao gồm lỗi nếu có)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logError("Lỗi thực thi scp:", err, "")
		return fmt.Errorf("SCP thất bại: %s - %s", err.Error(), string(output))
	}
	return nil
}

// Hàm xử lý tải file từ client và chuyển lên server từ xa bằng SCP
func TransferFile(c *gin.Context) {
	// Lấy file từ request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc file"})
		return
	}
	defer file.Close() // Đảm bảo file sẽ được đóng sau khi xử lý xong
	// Đảm bảo thư mục uploads tồn tại
	if _, err := os.Stat(config.AppConfig.Storage.UploadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(config.AppConfig.Storage.UploadPath, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo thư mục uploads"})
			return
		}
	}
	// Xây dựng đường dẫn lưu file cục bộ (local)
	localFilePath := filepath.Join(config.AppConfig.Storage.UploadPath, header.Filename)

	// Tạo file cục bộ để lưu trữ tạm thời
	outFile, err := os.Create(localFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo file tạm"})
		return
	}

	// Sao chép nội dung từ file được tải lên vào file cục bộ
	_, err = io.Copy(outFile, file)
	outFile.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":     "Lỗi khi lưu file",
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	// Xây dựng đường dẫn file trên server từ xa (remote)
	remoteFilePath := filepath.Join(config.AppConfig.SCP.RemotePath, header.Filename)

	// Gọi hàm SCP để chuyển file lên máy chủ từ xa
	err = scpFile(localFilePath, remoteFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Xóa file tạm sau khi SCP thành công
	if err := os.Remove(localFilePath); err != nil {
		logError("Lỗi khi xóa file tạm:", err, "")
	}

	// Trả về kết quả nếu chuyển file thành công
	c.JSON(http.StatusOK, gin.H{
		"message":     "File đã chuyển thành công",
		"remote_path": remoteFilePath,
		"status_code": http.StatusOK,
	})

}
