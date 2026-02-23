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

// Получаем данные из БД сохраняем в слайс на основе структуры
func GetWords(catId int) []Word {
	db.DataBase()
	defer db.DB.Close()
	var words []Word
	//READ
	rows, err := db.DB.Query("SELECT * FROM wordsToGuess WHERE catid = $1", catId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var w Word
		if err := rows.Scan(&w.Wordindex, &w.Hint, &w.Answer, &w.Categroy, &w.CatId); err != nil {
			log.Fatal(err)
		}
		words = append(words, w)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return words

}
