package data

import (
	"log"
	"wordGame/internal/models"
	db "wordGame/pkg"
)

func GetWords(catId int) []models.Word {
	database := db.GetDB()
	if database == nil {
		log.Fatal("DB is not initialized")
	}
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

	var words []models.Word

	for rows.Next() {
		var w models.Word
		if err := rows.Scan(&w.Wordindex, &w.Hint, &w.Answer, &w.Categroy, &w.CatId); err != nil {
			log.Fatal(err)
		}
		words = append(words, w)
	}

	return words
}

func GetCategories() []models.Categroy {
	database := db.GetDB()
	if database == nil {
		log.Fatal("DB is not initialized")
	}

	rows, err := database.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var categories []models.Categroy

	for rows.Next() {
		var c models.Categroy
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			log.Fatal(err)
		}
		categories = append(categories, c)
	}

	return categories
}
