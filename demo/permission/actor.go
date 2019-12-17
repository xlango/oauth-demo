package main

import "oauth-demo/demo/common"

func init() {
	common.CreateTalbe(Actor{})
}

type Actor struct {
	id   int64 `gorm:"AUTO_INCREMENT;primary_key"`
	name string
}
