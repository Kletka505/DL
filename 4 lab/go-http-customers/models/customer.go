package models

type Customer struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	UserPass string `gorm:"not null" json:"user_pass"`
}
