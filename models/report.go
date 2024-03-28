package model

import "time"

type Report struct {
	ID                  uint        `gorm:"primary_key;auto_increment" json:"id"`
	ReportName          string      `gorm:"varchar(255);not null" json:"reportName"`
	ReportGroupID       uint        `gorm:"not null" json:"reportGroupId"` // Foreign key referencing ReportGroup
	ReportGroup         ReportGroup `gorm:"foreignkey:ReportGroupID"`      // Define relationship with ReportGroup
	ReportDescripton    string      `gorm:"type:text;not null" json:"reportDescription"`
	ReportPriceCurrency string      `gorm:"varchar(255);default:BDT;not null" json:"reportPriceCurrency"`
	ReportPrice         string      `gorm:"varchar(255);not null" json:"reportPrice"`
	CreatedBy           string      `gorm:"not null" json:"createdBy"`
	CreatedAt           time.Time   `gorm:"not null" json:"createdAt"`
	UpdatedBy           string      `gorm:"not null" json:"updatedBy"`
	UpdatedAt           time.Time   `gorm:"not null" json:"updatedAt"`
	Mode                string      `gorm:"default:n" json:"mode"`
	Status              int64       `gorm:"default:1" json:"status"`
}
