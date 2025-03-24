package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongoDB() {
	uri := fmt.Sprintf("mongodb+srv://%s:%s@filesystem.gcqxk.mongodb.net/?retryWrites=true&w=majority&appName=%s",
		AppConfig.Database.Username, AppConfig.Database.Password, AppConfig.Database.Project)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Không thể kết nối MongoDB Atlas:", err)
	}
	fmt.Println("Đã kết nối MongoDB Atlas!")
	MongoClient = client
}

// Lấy collection user
func GetUserCollection() *mongo.Collection {
	return MongoClient.Database(AppConfig.Database.DbName).Collection(AppConfig.Database.Collection)
}
