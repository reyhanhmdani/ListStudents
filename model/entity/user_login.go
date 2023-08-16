package entity

type Admin struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	UserID   int64  `json:"-"`
	Username string `gorm:"type:varchar(255)" json:"username"`
	Email    string `gorm:"type:varchar(255);uniqueIndex:idx_email_name" json:"email"`
	Password string `json:"password" binding:"required"`
}

type UserLogin struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TableName specifies the custom table name for Mahasiswa entity
