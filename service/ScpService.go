package service

import (
	"fmt"
	"io"

	// "log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"backend/config"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// Hàm thực hiện lệnh SCP để chuyển file lên máy chủ từ xa
func scpFile(localFile, remoteFile string, remotePort string) error {
	cmd := exec.Command("scp", "-P", remotePort, localFile, // Tạo lệnh SCP
		fmt.Sprintf("%s@%s:%s", config.AppConfig.SCP.Username, config.AppConfig.SCP.RemoteHost, remoteFile))

	// Thực thi lệnh SCP và lấy đầu ra (bao gồm lỗi nếu có)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// logError("Lỗi thực thi scp:", err, "")
		return fmt.Errorf("SCP thất bại: %s - %s", err.Error(), string(output))
	}
	return nil
}

// Hàm xử lý tải file từ client và chuyển lên server từ xa bằng SCP
func TransferFile(c *gin.Context) {
	// Lấy file từ request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc file"})
		// utils.ErrorResponse("Không thể đọc file", http.StatusBadRequest, nil)
		return
	}
	defer file.Close() // Đảm bảo file sẽ được đóng sau khi xử lý xong

	remotePort := c.PostForm("port")
	remotePath := c.PostForm("path")
	if remotePort == "" {
		remotePort = config.AppConfig.SCP.RemotePort
	}

	ok, err := utils.CheckFolder(config.AppConfig.Storage.UploadPath); !ok {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể tạo thư mục uploads: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}
	// Xây dựng đường dẫn lưu file cục bộ (local)
	localFilePath := filepath.Join(config.AppConfig.Storage.UploadPath, header.Filename)

	// Tạo file cục bộ để lưu trữ tạm thời
	outFile, err := os.Create(localFilePath)
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Không thể tạo file tạm", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	// Sao chép nội dung từ file được tải lên vào file cục bộ
	_, err = io.Copy(outFile, file)
	outFile.Close()
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi khi lưu file", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	// Xây dựng đường dẫn file trên server từ xa (remote)
	remoteFilePath := filepath.Join(remotePath, header.Filename)
	fmt.Println(remoteFilePath)
	// Gọi hàm SCP để chuyển file lên máy chủ từ xa
	err = scpFile(localFilePath, remoteFilePath, remotePort)
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("SCP thất bại: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}
	// Xóa file tạm sau khi SCP thành công
	if err := os.Remove(localFilePath); err != nil {
		// logError("Lỗi khi xóa file tạm:", err, "")
	}
	utils.SuccessResponse(c, "File đã chuyển thành công", http.StatusOK, gin.H{"remote_path": remoteFilePath})

}
