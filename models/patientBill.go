package model

import "time"

type PatientBill struct {
	ID                    uint      `gorm:"primary_key;auto_increment" json:"id"`
	PatientID             uint      `gorm:"not null" json:"PatientId"`
	Patient               Patient   `gorm:"foreignkey:PatientID"`
	UniqueBillId		string `gorm:"not null" json:"uniqueBillId"`
	Reference             string    `gorm:"null" json:"reference"`
	AmountToPay           float64   `gorm:"not null" json:"amountToPay"`
	Discount              float32   `gorm:"null" json:"discount"`
	DiscountedAmountToPay float64   `gorm:"not null" json:"discountedAmountToPay"`
	AmountPaidReception   float64   `gorm:"not null" json:"amountPaidReception"`
	AmountPaidFinalStage  float64   `gorm:"not null" json:"amountPaidFinalStage"`
	AmountDue             float64   `gorm:"not null" json:"amountDue"`
	CreatedBy             string    `gorm:"not null" json:"createdBy"`
	CreatedAt             time.Time `gorm:"not null" json:"createdAt"`
	UpdatedBy             string    `gorm:"not null" json:"updatedBy"`
	UpdatedAt             time.Time `gorm:"not null" json:"updatedAt"`
	Mode                  string    `gorm:"default:n" json:"mode"`
	Status                int64     `gorm:"default:1" json:"status"`
}
