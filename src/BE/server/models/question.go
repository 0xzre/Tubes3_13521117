package models

type Question struct {
	Question *string `json:"question"`
	Answer   *string `json:"answer"`
}
