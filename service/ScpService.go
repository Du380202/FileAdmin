package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"backend/config"

	"github.com/gin-gonic/gin"
)

func scpFile(localFile, remoteFile string) error {
	cmd := exec.Command("scp", "-P", config.AppConfig.SCP.RemotePort, localFile,
		fmt.Sprintf("%s@%s:%s", config.AppConfig.SCP.Username, config.AppConfig.SCP.RemoteHost, remoteFile))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("SCP thất bại: %s - %s", err.Error(), string(output))
	}

	fmt.Println("SCP thành công:", string(output))
	return nil
}

func TransferFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Không thể đọc file"})
		return
	}
	defer file.Close()

	localFilePath := filepath.Join("uploads", header.Filename)
	outFile, err := os.Create(localFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Không thể tạo file tạm"})
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":     "Lỗi khi lưu file",
			"status_code": http.StatusInternalServerError,
		})
		return
	}

	remoteFilePath := filepath.Join(config.AppConfig.SCP.RemotePath, header.Filename)
	err = scpFile(localFilePath, remoteFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "File đã chuyển thành công",
		"remote_path": remoteFilePath,
		"status_code": http.StatusOK,
	})
}
