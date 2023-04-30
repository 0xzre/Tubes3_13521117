package models

type History struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	Prompt  string `gorm:"type:text;primaryKey" json:"prompt"`
	Respons string `gorm:"type:text;primaryKey" json:"respons"`
}
