package model

import "time"

type ReportGroup struct {
    ID                  uint      `gorm:"primary_key;auto_increment" json:"id"`
    GroupName           string    `gorm:"varchar(255);not null" json:"reportGroupName"`
    GroupDescripton     string    `gorm:"type:text;null" json:"reportGroupDescription"`
    CreatedBy           string    `gorm:"not null" json:"createdBy"`
    CreatedAt           time.Time `gorm:"not null" json:"createdAt"`
    UpdatedBy           string    `gorm:"not null" json:"updatedBy"`
    UpdatedAt           time.Time `gorm:"not null" json:"updatedAt"`
    Mode                string    `gorm:"default:n" json:"mode"`
    Status              int64     `gorm:"default:1" json:"status"`
}