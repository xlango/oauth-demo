package main

func init() {
	CreateTalbe(Actor{})
}

type Actor struct {
	id		int64 	`gorm:"AUTO_INCREMENT;primary_key"`
	name 	string
}