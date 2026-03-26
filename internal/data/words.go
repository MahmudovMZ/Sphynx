package data

import (
	"log"
	db "wordGame/pkg"
)

// Структура для подсказок и ответов номера загадки и категории слов с номером

type Word struct {
	Wordindex int
	Hint      string
	Answer    string
	Categroy  string
	CatId     int
}

type Categroy struct {
	ID   int    `json:"category_id"`
	Name string `json:"category_name"`
}

func GetWords(catId int) []Word {
	database := db.GetDB()
	if database == nil {
		log.Fatal("DB is not initialized")
	}

	rows, err := database.Query(
		"SELECT index, hint, answer, category, catid FROM wordsToGuess WHERE catid = $1",
		catId,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var words []Word

	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.Wordindex, &w.Hint, &w.Answer, &w.Categroy, &w.CatId); err != nil {
			log.Fatal(err)
		}
		words = append(words, w)
	}

	return words
}

func GetCategories() []Categroy {
	database := db.GetDB()
	if database == nil {
		log.Fatal("DB is not initialized")
	}

	rows, err := database.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []Categroy

	for rows.Next() {
		var c Categroy
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, c)
	}

	return categories
}
