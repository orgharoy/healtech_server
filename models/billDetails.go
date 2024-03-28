package model

import "time"

type BillDetails struct {
	ID            uint        `gorm:"primary_key;auto_increment" json:"id"`
	PatientBillID uint        `gorm:"not null" json:"PatientBillId"`
	PatientBill   PatientBill `gorm:"foreignkey:PatientBillID"`
	ReportID      uint        `gorm:"not null" json:"ReportId"`
	Report        Report      `gorm:"foreignkey:ReportID"`
	CreatedBy     string      `gorm:"not null" json:"createdBy"`
	CreatedAt     time.Time   `gorm:"not null" json:"createdAt"`
	UpdatedBy     string      `gorm:"not null" json:"updatedBy"`
	UpdatedAt     time.Time   `gorm:"not null" json:"updatedAt"`
	Mode          string      `gorm:"default:n" json:"mode"`
	Status        int64       `gorm:"default:1" json:"status"`
}
