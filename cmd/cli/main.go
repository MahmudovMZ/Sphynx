package main

import (
	ui "wordGame/internal/UI"
	Game "wordGame/internal/game"
)

func main() {
	//Запуск всей программы
	ui.PrintMsg()
	category := ui.ScanCat()
	ui.CategoryConfirm(category)
	game := Game.NewGame(category)
	game.Run()
}

// подключить БД вместо структуры --- DONE
// Вынести логику в Веб-запросы
