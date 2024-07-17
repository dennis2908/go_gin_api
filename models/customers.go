package models

// "time"

type Customer struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"type:varchar(300)" json:"name"`
	UserName string `gorm:"type:varchar(300)" json:"username"`
	Email    string `gorm:"type:text" json:"email"`
	Password string `gorm:"type:text" json:"password"`
}
