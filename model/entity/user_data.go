package entity

type User_data struct {
	ID          int64        `gorm:"primaryKey" json:"id"`
	Email       string       `gorm:"type:varchar(255);uniqueIndex:idx_email_name" json:"email"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex:idx_email_name" json:"name"`
	Age         int          `json:"age"`
	Address     string       `gorm:"type:varchar(100)" json:"address"`
	Birthdate   string       `gorm:"type:varchar(255)" json:"birthdate"`
	PhoneNumber string       `gorm:"type:varchar(20)" json:"phone_number"`
	UserID      int64        `json:"user_id"`
	Attachments []Attachment `gorm:"foreignKey:user_id" json:"attachments"`
}
