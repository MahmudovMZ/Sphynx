package models

type Word struct {
	Wordindex int
	Hint      string
	Answer    string
	Categroy  string
	CatId     int
}

type Category struct {
	ID   int
	Name string
}

type LeaderBoard struct {
	UserId   int
	UserName string
	Score    int
}
