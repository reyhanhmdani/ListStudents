package entity

type Roles struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	RoleName string `gorm:"not null" json:"role_name"`
}
