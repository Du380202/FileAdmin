// package config

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"time"

// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// var MongoClient *mongo.Client

// func ConnectMongoDB() {
// 	uri := fmt.Sprintf("mongodb+srv://%s:%s@filesystem.gcqxk.mongodb.net/?retryWrites=true&w=majority&appName=%s",
// 		AppConfig.Database.Username, AppConfig.Database.Password, AppConfig.Database.Project)

// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Kiểm tra kết nối
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		log.Fatal("Không thể kết nối MongoDB Atlas:", err)
// 	}
// 	fmt.Println("Đã kết nối MongoDB Atlas!")
// 	MongoClient = client
// }

// // Lấy collection user
// func GetUserCollection() *mongo.Collection {
// 	return MongoClient.Database(AppConfig.Database.DbName).Collection(AppConfig.Database.Collection)
// }

package config

import (
	"backend/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectMySQL() {
	// Chuỗi kết nối MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Database.Username,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.DbName,
	)

	// Kết nối với MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối MySQL:", err)
	}

	// Chạy AutoMigrate để tạo bảng nếu chưa có
	err = db.AutoMigrate(&models.User{}, &models.File{})
	if err != nil {
		log.Fatal("Lỗi khi chạy AutoMigrate:", err)
	}

	fmt.Println("Đã kết nối MySQL!")
	DB = db
}

// Lấy DB instance
func GetDB() *gorm.DB {
	return DB
}
