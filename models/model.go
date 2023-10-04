package models

import "time"

type BaseRepo struct {
	ID        uint      `gorm:"id;primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"create_at"`
	UpdatedAt time.Time `json:"update_at" gorm:"update_at"`
}
