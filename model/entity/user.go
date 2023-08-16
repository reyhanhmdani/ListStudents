package entity

type User struct {
	ID       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"not null;unique" json:"email"`
	Role     string `json:"role"`
}
