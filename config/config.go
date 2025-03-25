package config

import (
	"log"

	"github.com/spf13/viper" // Thư viện hỗ trợ đọc và quản lý file cấu hình
)

// Định nghĩa cấu trúc Config chứa thông tin cấu hình cho server và storage
type Config struct {
	Server   ServerConfig  // Cấu hình server (port)
	Storage  StorageConfig // Cấu hình lưu trữ (đường dẫn upload)
	SCP      SCPConfig
	Database DatabaseConfig
}

// Cấu trúc cấu hình cho server
type ServerConfig struct {
	Port string // Cổng chạy server
}

// Cấu trúc cấu hình cho lưu trữ
type StorageConfig struct {
	UploadPath string // Đường dẫn lưu trữ file upload
}

type SCPConfig struct {
	RemoteHost string // Địa chỉ máy chủ từ xa
	RemotePort string // Cổng SSH (thường là 22)
	Username   string // Tên đăng nhập SSH
	RemotePath string // Đường dẫn thư mục trên máy chủ đích
}

type DatabaseConfig struct {
	Project    string
	Username   string
	Password   string
	DbName     string
	Collection string
}

// Biến toàn cục lưu trữ cấu hình của ứng dụng
var AppConfig Config

// LoadConfig đọc file config và gán giá trị vào biến AppConfig
func LoadConfig() {
	viper.SetConfigName("config") // Đặt tên file config (không cần phần mở rộng)
	viper.SetConfigType("yaml")   // Chỉ định kiểu file config là YAML
	viper.AddConfigPath("config") // Đặt đường dẫn thư mục chứa file config

	// Đọc file config, nếu lỗi thì dừng chương trình
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Lỗi đọc config file: %v", err)
	}

	// Parse dữ liệu từ file config vào biến AppConfig
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
