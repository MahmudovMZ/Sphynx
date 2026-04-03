package models

type Bot_Menu struct {
	Id    int
	Title string
}

var Menu = []Bot_Menu{
	{Id: 1, Title: "Registration"},
	{Id: 2, Title: "Login"},
	{Id: 3, Title: "Leader Board"},
	{Id: 4, Title: "Quit"},
}
