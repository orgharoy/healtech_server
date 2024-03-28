package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	UserName  string    `gorm:"varchar(255);not null" json:"userName"`
	UserId    string    `gorm:"varchar(255);not null" json:"userId"`
	Password  string    `gorm:"not null" json:"password"`
	Role      string    `gorm:"not null" json:"role"`
	CreatedBy string    `gorm:"not null" json:"createdBy"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt"`
	UpdatedBy string    `gorm:"not null" json:"updatedBy"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt"`
	Mode      string    `gorm:"default:n" json:"mode"`
	Status    int64     `gorm:"default:1" json:"status"`
}
