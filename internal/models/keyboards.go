package models

type Bot_Menu struct {
	Id    int
	Title string
}

var BotMenu = []Bot_Menu{
	{Id: 1, Title: "Registration"},
	{Id: 2, Title: "Login"},
}

var GameMenu = []Bot_Menu{
	{Id: 1, Title: "New Game"},
	{Id: 2, Title: "Leader Board"},
	{Id: 3, Title: "Quit"},
	{Id: 3, Title: "Exit"},
}
