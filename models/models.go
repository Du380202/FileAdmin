package models

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Folder   string `gorm:"not null" json:"folder"`
	Files    []File `gorm:"foreignKey:UserID" json:"files"` // Liên kết với File
}

type File struct {
	FileId   uint   `gorm:"primaryKey;autoIncrement" json:"fileid"`
	FileName string `gorm:"unique;not null" json:"filename"`
	Content  string `gorm:"not null" json:"content"`
	FilePath string `gorm:"not null" json:"filepath"`
	UserID   uint   `gorm:"not null" json:"userid"`                                    // Khóa ngoại
	// User     User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"` // Ràng buộc xóa
}
