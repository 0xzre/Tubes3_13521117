package models

type Question struct {
	Pertanyaan string `gorm:"type:text;primaryKey" json:"pertanyaan"`
	Jawaban    string `gorm:"type:text" json:"jawaban"`
}
