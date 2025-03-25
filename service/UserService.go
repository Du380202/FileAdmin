package service

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// API đăng ký user
func Register(c *gin.Context) {
	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("error %s", err.Error()), http.StatusBadRequest, nil)

		return
	}

	// Hash mật khẩu
	hashedPassword, err := middleware.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi mã hóa mật khẩu: %s", err.Error()), http.StatusInternalServerError, nil)

		return
	}
	req.Password = hashedPassword

	// Thêm user vào MongoDB
	userCollection := config.GetUserCollection()
	_, err = userCollection.InsertOne(context.TODO(), req)
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi khi thêm user: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	utils.SuccessResponse(c, "Đăng ký thành công", http.StatusCreated, nil)
}

// API đăng nhập user
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi đăng nhập: %s", err.Error()), http.StatusBadRequest, nil)
		return
	}

	// Lấy user từ MongoDB
	userCollection := config.GetUserCollection()
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		utils.ErrorResponse(c, "Sai tài khoản hoặc mật khẩu", http.StatusUnauthorized, nil)
		return
	}

	// Kiểm tra mật khẩu
	if !middleware.CheckPasswordHash(req.Password, user.Password) {
		utils.ErrorResponse(c, "Sai tài khoản hoặc mật khẩu", http.StatusUnauthorized, nil)
		return
	}

	// Tạo JWT Token
	token, err := middleware.GenerateJWT(user.Username)
	if err != nil {
		utils.ErrorResponse(c, fmt.Sprintf("Lỗi tạo token: %s", err.Error()), http.StatusInternalServerError, nil)
		return
	}

	utils.SuccessResponse(c, "Đăng ký thành công", http.StatusOK, gin.H{"token": token})
}
