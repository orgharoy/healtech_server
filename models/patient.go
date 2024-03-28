package model

import "time"

type Patient struct {
	ID              uint      `gorm:"primary_key;auto_increment" json:"id"`
	UniquePatientId string    `gorm:"not null;unique_index" json:"uniquePatientId"`
	Name            string    `gorm:"not null" json:"patientName"`
	Age             string    `gorm:"not null" json:"patientAge"`
	Gender          string    `gorm:"not null" json:"patientGender"`
	ContactNumber   string    `gorm:"not null" json:"patientContactNumber"`
	CreatedBy       string    `gorm:"not null" json:"createdBy"`
	CreatedAt       time.Time `gorm:"not null" json:"createdAt"`
	UpdatedBy       string    `gorm:"not null" json:"updatedBy"`
	UpdatedAt       time.Time `gorm:"not null" json:"updatedAt"`
	Mode            string    `gorm:"default:n" json:"mode"`
	Status          int64     `gorm:"default:1" json:"status"`
}
