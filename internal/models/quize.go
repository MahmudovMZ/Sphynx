package models

type Word struct {
	Wordindex int
	Hint      string
	Answer    string
	Categroy  string
	CatId     int
}

type Category struct {
	ID   int    `json:"category_id"`
	Name string `json:"category_name"`
}
